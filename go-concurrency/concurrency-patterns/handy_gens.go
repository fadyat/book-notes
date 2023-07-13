package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// will generate the same values until not stopped
	repeat := func(
		done <-chan interface{},
		values ...interface{},
	) <-chan interface{} {
		stream := make(chan interface{})

		go func() {
			defer close(stream)

			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case stream <- v:
					}
				}
			}
		}()

		return stream
	}

	// will take the first num values from the stream
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

	// will generate values from fn until not stopped
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

	done := make(chan interface{})
	defer close(done)

	for num := range take(done, repeat(done, 1, 2, 3), 10) {
		fmt.Println(num)
	}

	rnd := func() interface{} { return rand.Int() }
	for num := range take(done, repeatFn(done, rnd), 10) {
		fmt.Println(num)
	}

	// will convert the stream to chan string
	toString := func(
		done <-chan interface{},
		stream <-chan interface{},
	) <-chan string {
		stringStream := make(chan string)

		go func() {
			defer close(stringStream)

			for v := range stream {
				select {
				case <-done:
					return
				case stringStream <- v.(string):
				}
			}
		}()

		return stringStream
	}

	for str := range toString(done, take(done, repeat(done, "a", "b", "c"), 10)) {
		fmt.Println(str)
	}
}
