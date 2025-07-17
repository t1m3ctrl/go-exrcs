package main

import (
	"fmt"
)

func SetBit(n int64, i int, bit int) int64 {
	if bit == 1 {
		return n | (1 << i)
	} else {
		return n &^ (1 << i)
	}
}

func main() {
	var num int64 = 31
	var result int64

	result = SetBit(num, 0, 0)
	fmt.Printf("Число %d (%b) после установки %d-го бита в %d: %d (%b)\n",
		num, num, 0, 0, result, result)

	result = SetBit(num, 10, 1)
	fmt.Printf("Число %d (%b) после установки %d-го бита в %d: %d (%b)\n",
		num, num, 10, 1, result, result)
}
