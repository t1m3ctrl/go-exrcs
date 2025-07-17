package main

import (
	"fmt"
	"math/rand/v2"
)

func intersect(arr1, arr2 []int) []int {
	mp := map[int]bool{}
	res := []int{}

	for _, v := range arr1 {
		mp[v] = true
	}

	for _, v := range arr2 {
		if mp[v] == true {
			res = append(res, v)
			mp[v] = false
		}
	}

	return res
}

func main() {
	arr1 := []int{}
	arr2 := []int{}

	n := 10

	for i := 0; i < n; i++ {
		num1 := rand.IntN(n)
		arr1 = append(arr1, num1)

		num2 := rand.IntN(n)
		arr2 = append(arr2, num2)
	}

	fmt.Printf("Array 1: %v\n", arr1)
	fmt.Printf("Array 2: %v\n", arr2)
	fmt.Printf("Intersection: %v\n", intersect(arr1, arr2))
}
