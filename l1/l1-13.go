package main

import (
	"fmt"
	"math/rand/v2"
)

func main() {
	a, b := rand.IntN(100), rand.IntN(100)
	fmt.Println(a, b)

	// a, b = b, a

	a = a - b
	b = b + a
	a = b - a

	// a = a ^ b
	// b = a ^ b
	// a = a ^ b

	fmt.Println(a, b)
}
