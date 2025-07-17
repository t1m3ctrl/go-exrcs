package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	notSafe()
	safe()
}

func notSafe() {
	var counter int
	var wg sync.WaitGroup

	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			counter++
		}()
	}

	wg.Wait()
	fmt.Println("Not safe: ", counter)
}

func safe() {
	var counter int32
	var wg sync.WaitGroup

	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			atomic.AddInt32(&counter, 1)
		}()
	}

	wg.Wait()
	fmt.Println("Safe: ", counter)
}
