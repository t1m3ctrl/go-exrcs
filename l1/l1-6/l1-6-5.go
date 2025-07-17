package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	// В моем понимании - та же остановка через канал
	timeout := time.After(5 * time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-timeout:
				fmt.Println("Остановлена по таймауту")
				return
			default:
				fmt.Println("Работает...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	wg.Wait()
}
