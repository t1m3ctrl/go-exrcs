package main

import (
	"fmt"
	"sort"
)

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

func main() {
	arr := []int{5, 3, 8, 4, 2, 7, 1, 10}
	fmt.Println("Изначальный массив:    ", arr)
	sorted := quickSort(arr)
	sort.Ints(arr)
	fmt.Println("Отсортированный массив:", sorted)
	fmt.Println("Go sort:               ", arr)
}
