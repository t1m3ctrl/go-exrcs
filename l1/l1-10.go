package main

import (
	"fmt"
)

func main() {
	// arr := []float64{}
	arr := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
	groups := map[int][]float64{}

	// for i := 0; i < 100; i++ {
	// 	num := math.Round(rand.Float64()*10)/10 + float64(rand.IntN(40))
	// 	if i%2 == 0 {
	// 		num *= -1
	// 	}
	// 	arr = append(arr, num)
	// 	// fmt.Println(arr[i])
	// }

	// sort.Float64s(arr)

	for i := range arr {
		group := int((arr[i] / 10)) * 10
		groups[group] = append(groups[group], arr[i])
	}

	for k, v := range groups {
		fmt.Printf("%d: %.1f\n", k, v)
	}
}
