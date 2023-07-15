package main

import (
	"fmt"
	"time"
)

func DoWorkLongLiving(
	done <-chan interface{},
	nums ...int,
) (<-chan interface{}, <-chan int) {
	// This ensures that there’s always at least one pulse
	// sent out even if no one is listening in time for to send to occur.
	heartbeatStream := make(chan interface{}, 1)
	workStream := make(chan int)

	go func() {
		defer func() {
			close(heartbeatStream)
			close(workStream)
		}()

		// Here we simulate some kind of delay before the goroutine can begin working.
		// In practice this can be all kinds of things and is nondeterministic.
		time.Sleep(2 * time.Second)

		for _, n := range nums {
			// We don’t want to include this in the same select
			// block as the send on results because if the receiver
			// isn’t ready for the result, they'll receive a pulse instead,
			// and the current value of the result will be lost.
			//
			// We also don’t include a case statement for the done channel
			// since we have a default case that will just fall through.
			select {
			case heartbeatStream <- struct{}{}:
			// Once again we guard against the fact that no one may be
			// listening to our heartbeats
			default:
			}

			select {
			case <-done:
				return
			case workStream <- n:
			}
		}
	}()

	return heartbeatStream, workStream
}

func DoWorkLongLivingLabel(
	done <-chan interface{},
	pulseInterval time.Duration,
	nums ...int,
) (<-chan interface{}, <-chan int) {
	heartbeatStream := make(chan interface{}, 1)
	workStream := make(chan int)

	go func() {
		defer func() {
			close(heartbeatStream)
			close(workStream)
		}()

		time.Sleep(2 * time.Second)

		pulse := time.Tick(pulseInterval)

		// we're using a label here to make continuing
		// from the inner loop a little simpler.
	numLoop:
		for _, n := range nums {
			for {
				select {
				case <-done:
					return
				case <-pulse:
					select {
					case heartbeatStream <- struct{}{}:
					default:
					}
				case workStream <- n:
					continue numLoop
				}
			}
		}
	}()

	return heartbeatStream, workStream
}

func main() {
	done := make(chan interface{}, 1)
	defer close(done)

	heartbeat, results := DoWorkLongLiving(done, 1, 2, 3)
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

			fmt.Printf("results %v\n", r)
		}
	}
}
