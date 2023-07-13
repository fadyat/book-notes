package main

import "fmt"

func main() {
	var orDone func(done <-chan interface{}, c <-chan interface{}) <-chan interface{}

	bridge := func(
		done <-chan interface{},
		chanStream <-chan <-chan interface{},
	) <-chan interface{} {
		// this channel will return all the values from the channels
		stream := make(chan interface{})

		go func() {
			defer close(stream)

			for {
				// pooling the channels
				var ch <-chan interface{}
				select {
				case maybeCh, ok := <-chanStream:
					if !ok {
						return
					}

					ch = maybeCh
				case <-done:
					return
				}

				// reading the values from the channel
				for val := range orDone(done, ch) {
					select {
					case stream <- val:
					case <-done:
					}
				}
			}
		}()

		return stream
	}

	genValues := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))

		go func() {
			defer close(chanStream)

			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()

		return chanStream
	}

	done := make(chan interface{})
	defer close(done)

	for val := range bridge(done, genValues()) {
		fmt.Println(val)
	}
}
