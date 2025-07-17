package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

func genNums(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(ch)

	arr := []int{}
	for i := 0; i < 10; i++ {
		arr = append(arr, rand.IntN(100))
	}

	for i := range arr {
		ch <- arr[i]
		fmt.Println("Wrote ", arr[i])
		// time.Sleep(300 * time.Millisecond)
	}
}

func squareNums(ch chan int, sqCh chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(sqCh)

	for {
		select {
		case data, ok := <-ch:
			if !ok {
				fmt.Println("Data ended in squareNums func from channel.")
				return
			}
			sqCh <- data * data
		}
	}
}

func printData(sqCh chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case data, ok := <-sqCh:
			if !ok {
				fmt.Println("Data ended in printData func from square channel.")
				return
			}
			fmt.Println("Squared ", data)
		}
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(3)

	ch := make(chan int)
	sqCh := make(chan int)

	go genNums(ch, &wg)
	go squareNums(ch, sqCh, &wg)
	go printData(sqCh, &wg)

	wg.Wait()
}
