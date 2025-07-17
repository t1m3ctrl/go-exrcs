package main

import (
	"fmt"
)

func bSearch(arr []int, target int) int {
	left := 0
	right := len(arr) - 1

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return -1
}

func main() {
	arr := []int{50, 30, 80, 40, 20, 70, 10, 100}
	sorted := quickSort(arr)
	fmt.Println(bSearch(sorted, 30))
	fmt.Println(bSearch(sorted, 128))
}

func quickSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	mid := arr[len(arr)/2]

	var less, equal, greater []int
	for _, value := range arr {
		switch {
		case value < mid:
			less = append(less, value)
		case value == mid:
			equal = append(equal, value)
		case value > mid:
			greater = append(greater, value)
		}
	}

	return append(append(quickSort(less), equal...), quickSort(greater)...)
}
