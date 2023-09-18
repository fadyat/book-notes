package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)
	defer close(done)

	go func() {
		select {
		case s := <-signalChan:
			fmt.Println("Got signal:", s)
			done <- true
		}
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}
