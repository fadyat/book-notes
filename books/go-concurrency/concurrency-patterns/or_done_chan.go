package main

import "fmt"

func main() {
	orDone := func(done, c <-chan interface{}) <-chan interface{} {
		stream := make(chan interface{})

		go func() {
			defer close(stream)

			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if !ok {
						return
					}

					select {
					case <-done:
					case stream <- v:
					}
				}
			}
		}()

		return stream
	}

	done := make(chan interface{})
	defer close(done)

	values := make(chan interface{})

	go func() {
		defer close(values)
		for i := 0; i < 10; i++ {
			values <- i
		}
	}()

	for val := range orDone(done, values) {
		fmt.Println(val)
	}
}
