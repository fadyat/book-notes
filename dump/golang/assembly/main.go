package main

import (
	"fmt"
	"math"
)

// go build will automatically detect sqrt_$GOARCH.s file
// and link it with the main.go file
func sqrt(x float64) float64

func main() {
	var number float64 = 2000

	fmt.Println(sqrt(number))
	fmt.Println(math.Sqrt(number))
}
