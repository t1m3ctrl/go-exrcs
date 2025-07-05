package main

import (
	"fmt"
	"sync"
)

func main() {
	nums := [5]int{2, 4, 6, 8, 10}
	wg := sync.WaitGroup{}

	wg.Add(len(nums))
	for _, v := range nums {
		go func(val int) {
			defer wg.Done()
			fmt.Println(val * val)
		}(v)
	}
	wg.Wait()

	wg.Add(len(nums))
	for _, v := range nums {
		val := v
		go func() {
			defer wg.Done()
			fmt.Println(val * val)
		}()
	}
	wg.Wait()
}
