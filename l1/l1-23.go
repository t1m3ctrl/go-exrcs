package main

import "fmt"

func main() {
	i := 5

	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	fmt.Println(arr)

	arr = append(arr[:i], arr[i+1:]...)
	fmt.Println(arr)

	newArr := make([]int, len(arr))
	copy(newArr, arr)
	fmt.Println(newArr)
}
