package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer fmt.Println("Завершение через Goexit")

		time.Sleep(500 * time.Millisecond)
		fmt.Println("Какой то код..")

		time.Sleep(500 * time.Millisecond)
		fmt.Println("Вызов Goexit()")
		runtime.Goexit() // Остановка горутины

		// Unreachable
		for {
			fmt.Println("Недостижимый код...")
			time.Sleep(100 * time.Millisecond)
		}
	}()

	wg.Wait()
}
