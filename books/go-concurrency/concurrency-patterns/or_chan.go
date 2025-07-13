package main

import (
	"fmt"
	"time"
)

func main() {
	// or function accepts a variadic number of channels and returns a single channel.
	var or func(channels ...<-chan interface{}) <-chan interface{}

	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		// this is recursive function -> stop point
		case 0:
			return nil
		// have only one channel -> return it
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})

		// we can wait for messages from multiple channels w/o blocking
		go func() {
			defer close(orDone)

			switch len(channels) {
			// an optimization, every recursive call will have at least two channels
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			// recursively create an or-channel from all channels and then select from it
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

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})

		go func() {
			defer close(c)
			time.Sleep(after)
		}()

		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	// we will see that the function returns after 1 second
	fmt.Printf("done after %v", time.Since(start))
}
