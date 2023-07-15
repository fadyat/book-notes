package main

import (
	"fmt"
	"time"
)

func main() {
	// Notice that because we might be sending out multiple pulses
	// while we wait for input, or multiple pulses while waiting to send results,
	// all the select statements need to be within for loops.

	doWork := func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) (<-chan interface{}, <-chan time.Time) {

		// set up a channel to send heartbeats on.
		heartbeat := make(chan interface{})
		results := make(chan time.Time)

		go func() {
			// not closing channels to simulate panic at goroutine
			//
			// defer func() {
			//	 close(heartbeat)
			//	 close(results)
			// }()

			// creating channel, which will make a pulse every interval
			pulse := time.Tick(pulseInterval)

			// another ticker to simulate work coming on.
			// interval picked higher so that we can see some heartbeats coming out
			// of goroutine
			workGen := time.Tick(2 * pulseInterval)

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				// including default clause, because we want guard against fact
				// that no one may be listening to our heartbeat.
				//
				// heartbeat results aren't critical
				default:
				}
			}

			sendResult := func(r time.Time) {
				for {
					select {
					case <-done:
						return
					case <-pulse:
						sendPulse()
					case results <- r:
						return
					}
				}
			}

			// Here is our simulated panic. Instead of infinitely looping until
			// we're asked to stop, as in the previous example, we'll only loop twice.
			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()

		return heartbeat, results
	}

	// standard done channel with 10 seconds of work
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)

	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}

			fmt.Println("pulse")
		case r, ok := <-results:
			if !ok {
				return
			}

			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			fmt.Println("worker goroutine isn't healthy!")
			return
		}
	}
}
