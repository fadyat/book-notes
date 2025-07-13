package main

import (
	"log"
	"os"
	"time"
)

func main() {
	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})

		go func() {
			defer close(orDone)

			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()

		return orDone
	}

	// Here we define the signature of a goroutine that can be monitored and
	// restarted.
	// We see the familiar done channel, and pulseInterval and heartbeat from
	// the heartbeat pattern.
	type startGoroutineFn func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) (heartbeat <-chan interface{})

	// On this line we see that a steward takes in a timeout for the goroutine
	// it will be monitoring, and a function, startGoroutine, to start the
	// goroutine it's monitoring.
	// Interestingly, the steward itself returns a startGoroutineFn indicating
	// that the steward itself is also monitorable.
	newSteward := func(
		timeout time.Duration,
		startGoroutine startGoroutineFn,
	) startGoroutineFn {
		return func(
			done <-chan interface{},
			pulseInterval time.Duration,
		) <-chan interface{} {
			heartbeat := make(chan interface{})

			go func() {
				defer close(heartbeat)

				var wardDone chan interface{}
				var wardHeartbeat <-chan interface{}

				// Here we define a closure that encodes a consistent
				// way to start the goroutine we're monitoring.
				startWard := func() {

					// This is where we create a new channel that we'll
					// pass into the ward goroutine in case we need to signal
					// that it should halt.
					wardDone = make(chan interface{})

					// Here we start the goroutine we'll be monitoring.
					// We want the ward goroutine to halt if either the steward
					// is halted, or the steward wants to halt the ward goroutine,
					// so we wrap both done channels in a logical-or.
					wardHeartbeat = startGoroutine(or(wardDone, done), timeout/2)
				}

				startWard()
				pulse := time.Tick(pulseInterval)

			monitorLoop:
				for {
					timeoutSignal := time.After(timeout)

					// This is our inner loop, which ensures that the steward
					// can send out pulses of its own.
					for {
						select {
						case <-pulse:
							select {
							case heartbeat <- struct{}{}:
							default:
							}
						// Here we see that if we receive the ward's pulse,
						// we continue our monitoring loop.
						case <-wardHeartbeat:
							continue monitorLoop
						// This line indicates that if we don't receive a pulse
						// from the ward within our timeout period, we request
						// that the ward halt, and we begin a new ward goroutine.
						case <-timeoutSignal:
							log.Printf("steward: ward unhealthy; restarting")
							close(wardDone)
							startWard()
							continue monitorLoop
						case <-done:
							return
						}
					}
				}
			}()

			return heartbeat
		}
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	// here we see that this goroutine isn’t doing anything but
	// waiting to be canceled. It’s also not sending out any pulses.
	doWork := func(
		done <-chan interface{},
		_ time.Duration,
	) <-chan interface{} {
		log.Printf("ward: hello, i'm irresponsible")
		go func() {
			<-done
			log.Printf("ward: i'm healing")
		}()

		return nil
	}

	doWorkWithSteward := newSteward(4*time.Second, doWork)

	done := make(chan interface{})
	time.AfterFunc(9*time.Second, func() {
		log.Printf("main: halting steward and ward")
		close(done)
	})

	for range doWorkWithSteward(done, 4*time.Second) {
	}
}
