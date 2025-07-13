package main

import "sync"

type Button struct {
	Clicked *sync.Cond
}

func main() {
	button := Button{
		sync.NewCond(&sync.Mutex{}),
	}

	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)

		go func() {
			goroutineRunning.Done()

			c.L.Lock()
			defer c.L.Unlock()

			c.Wait()
			fn()
		}()

		goroutineRunning.Wait()
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)

	subscribe(button.Clicked, func() {
		println("Maximizing window.")
		clickRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		println("Mouse clicked.")
		clickRegistered.Done()
	})

	// Broadcast wakes all goroutines waiting on c.
	button.Clicked.Broadcast()

	clickRegistered.Wait()
}
