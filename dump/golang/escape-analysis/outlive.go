package main

import (
	"fmt"
	"math/rand"
)

func getRandomPtr() *int {
	num := rand.Intn(100)
	return &num
}

func functionReturn() {
	// Any returned value outlives the function since
	// the called function does not know the value.

	num := getRandomPtr()
	fmt.Println(*num)
}

func loop() {
	// Variables declared outside a loop outlive the
	// assignment within the loop:

	var ptr *int
	for i := 0; i < 100; i++ {
		ptr = new(int)
		*ptr = i
	}

	fmt.Println(*ptr)
}

func clojure() {
	// Variables declared outside a closure outlive
	// the assignment within the closure

	var ptr *int
	func() {
		ptr = new(int)
		*ptr = 10
	}()

	fmt.Println(*ptr, ptr)
}

func insufficientStackSpace() {
	size := 1000000
	s := make([]int, size, size)
	for i := 0; i < size; i++ {
		s[i] = i
	}

	fmt.Println(len(s))
}

func dynamicTypeEscape() {
	fmt.Println("dynamicTypeEscape")
}

func main() {
	functionReturn()
	loop()
	clojure()
	insufficientStackSpace()
	dynamicTypeEscape()
}
