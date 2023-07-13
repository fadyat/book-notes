package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	// create a new condition variable
	c := sync.NewCond(&sync.Mutex{})

	// queue is empty at this point, later we add some entries
	queue := make([]interface{}, 0, 10)

	// wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	remove := func(delay time.Duration) {
		defer wg.Done()
		time.Sleep(delay)

		// entering the critical section, lock the mutex
		c.L.Lock()
		queue = queue[1:]
		log.Println("removed from queue: ", len(queue))

		// leaving the critical section, unlock the mutex
		c.L.Unlock()

		// signal the condition that the queue changed
		//
		// notifies goroutines, that are waiting for the condition
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)

		// entering the critical section, lock the mutex
		c.L.Lock()

		// waiting for the condition to be true
		for len(queue) == 2 {

			// suspends the execution of the calling goroutine
			c.Wait()
		}

		log.Println("adding to queue: ", len(queue))
		queue = append(queue, struct{}{})

		// calling new goroutine to remove from queue
		go remove(1 * time.Second)

		// leaving the critical section, unlock the mutex
		c.L.Unlock()
	}

	wg.Wait()
}
