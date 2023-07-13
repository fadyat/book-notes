package main

func main() {
	owner := func() <-chan int {
		ch := make(chan int, 5)

		go func() {
			defer close(ch)

			for i := 0; i < 5; i++ {
				ch <- i
			}
		}()

		return ch
	}

	results := owner()
	for result := range results {
		println("result:", result)
	}

	println("done")
}
