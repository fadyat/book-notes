package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	// fanIn is a function that takes a variadic number of channels
	// 	and multiplexes them onto a single channel.
	fanIn := func(
		done <-chan interface{},
		channels ...<-chan interface{},
	) <-chan interface{} {
		var wg sync.WaitGroup
		multiplexedStream := make(chan interface{})

		multiplex := func(c <-chan interface{}) {
			defer wg.Done()

			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}

		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}

	repeatFn := func(
		done <-chan interface{},
		fn func() interface{},
	) <-chan interface{} {
		stream := make(chan interface{})

		go func() {
			defer close(stream)

			for {
				select {
				case <-done:
					return
				case stream <- fn():
				}
			}
		}()

		return stream
	}

	take := func(
		done <-chan interface{},
		stream <-chan interface{},
		num int,
	) <-chan interface{} {
		takeStream := make(chan interface{})

		go func() {
			defer close(takeStream)

			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-stream:
				}
			}
		}()

		return takeStream
	}

	isPrime := func(n int) bool {
		if n < 2 {
			return false
		}

		for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
			if n%i == 0 {
				return false
			}
		}

		return true
	}

	primeFinder := func(
		done <-chan interface{},
		stream <-chan interface{},
	) <-chan interface{} {
		primeStream := make(chan interface{})

		go func() {
			defer close(primeStream)

			for {
				select {
				case <-done:
					return
				case i := <-stream:
					if isPrime(i.(int)) {
						primeStream <- i
					}
				}
			}
		}()

		return primeStream
	}

	done := make(chan interface{})
	defer close(done)

	rnd := func() interface{} { return rand.Intn(5_000_000) }
	stream := repeatFn(done, rnd)

	findersNumber := runtime.NumCPU()
	finders := make([]<-chan interface{}, findersNumber)
	fmt.Printf("Spinning up %d prime finders.\n", findersNumber)
	for i := 0; i < findersNumber; i++ {
		finders[i] = primeFinder(done, stream)
	}

	start := time.Now()
	fmt.Println("Primes:")
	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v\n", time.Since(start))
}
