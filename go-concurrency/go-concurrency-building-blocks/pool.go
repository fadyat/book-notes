package main

import (
	"fmt"
	"sync"
)

func main() {
	var numCalcsCreated int
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated++
			mem := make([]byte, 1024)
			return &mem // pointer to a slice
		},
	}

	// seeding the pool with 4KB
	for i := 0; i < 4; i++ {
		calcPool.Put(calcPool.New())
	}

	const numWorkers = 1024 * 1024

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			mem := calcPool.Get().(*[]byte) // type assertion
			defer calcPool.Put(mem)

			// do something interesting with mem
		}()
	}

	wg.Wait()
	fmt.Printf("%d calculators were created.", numCalcsCreated)
}
