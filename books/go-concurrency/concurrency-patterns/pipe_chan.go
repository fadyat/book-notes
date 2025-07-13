package main

import "fmt"

func main() {

	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		stream := make(chan int)

		go func() {
			defer close(stream)
			for _, i := range integers {
				select {
				case <-done:
					return
				case stream <- i:
				}
			}
		}()

		return stream
	}

	multiply := func(
		done <-chan interface{},
		stream <-chan int,
		multiplier int,
	) <-chan int {
		multipliedStream := make(chan int)

		go func() {
			defer close(multipliedStream)

			for i := range stream {
				select {
				case <-done:
					return
				case multipliedStream <- i * multiplier:
				}
			}
		}()

		return multipliedStream
	}

	add := func(
		done <-chan interface{},
		stream <-chan int,
		additive int,
	) <-chan int {
		addedStream := make(chan int)

		go func() {
			defer close(addedStream)

			for i := range stream {
				select {
				case <-done:
					return
				case addedStream <- i + additive:
				}
			}
		}()

		return addedStream
	}

	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}
