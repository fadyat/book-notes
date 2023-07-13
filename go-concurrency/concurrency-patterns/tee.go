package main

import "fmt"

func main() {
	var repeat func(done <-chan interface{}, values ...interface{}) <-chan interface{}
	var take func(done <-chan interface{}, stream <-chan interface{}, num int) <-chan interface{}
	var orDone func(done, c <-chan interface{}) <-chan interface{}

	tee := func(
		done <-chan interface{},
		in <-chan interface{},
	) (_, _ <-chan interface{}) {
		out1 := make(chan interface{})
		out2 := make(chan interface{})

		go func() {
			defer func() {
				close(out1)
				close(out2)
			}()

			for val := range orDone(done, in) {
				// we want to use local versions of out1 and out2,
				//  so we shadow them
				var out1, out2 = out1, out2

				// making for loop to send val to out1 and out2
				//  don't block each other
				for i := 0; i < 2; i++ {
					// once we read from out1 or out2, we set shadowed
					//  copies to nil, so that we don't send to nil channels
					select {
					case <-done:
					case out1 <- val:
						out1 = nil
					case out2 <- val:
						out2 = nil
					}
				}
			}
		}()

		return out1, out2
	}

	done := make(chan interface{})
	defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 1, 2, 3, 4, 5), 10))
	for val1 := range out1 {
		fmt.Printf("out1: %v\n", val1)
		fmt.Printf("out2: %v\n", <-out2)
	}
}
