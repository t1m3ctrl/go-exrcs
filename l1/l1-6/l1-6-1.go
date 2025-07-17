package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	stopFlag := false

	wg.Add(1)
	go func() {
		defer wg.Done()
		for !stopFlag {
			time.Sleep(500 * time.Millisecond)
			fmt.Println("Работает...")
		}
		fmt.Println("Остановка по условию")
	}()

	time.Sleep(5 * time.Second)

	stopFlag = true // Условие для остановки
	wg.Wait()
}
