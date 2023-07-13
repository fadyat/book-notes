package main

import (
	"fmt"
	"time"
)

func putSomeData(stream chan<- interface{}) {
	defer func() {
		close(stream)
		fmt.Println("close stream")
	}()

	for i := 0; i < 5; i++ {
		stream <- i
	}
}

func main() {
	// making buffered channel to close channel before reading
	stream := make(chan interface{}, 5)

	// blocking read = deadlock
	//  `ok` returned only if channel is closed or have some data
	//
	// v, ok := <-stream
	// fmt.Printf("%v, %v\n", v, ok)

	go putSomeData(stream)

	time.Sleep(time.Second)
	fmt.Println("sleep done")
	for i := 0; i < 10; i++ {
		v, ok := <-stream
		fmt.Printf("%v, %v\n", v, ok)
	}

	stream = make(chan interface{}, 5)
	go putSomeData(stream)

	// range will perform auto exit if channel is closed
	for v := range stream {
		fmt.Printf("%v\n", v)
	}
}
