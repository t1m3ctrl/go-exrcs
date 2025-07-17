package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	stopCh := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopCh:
				fmt.Println("Остановка через канал")
				return
			default:
				time.Sleep(500 * time.Millisecond)
				fmt.Println("Работает...")
			}
		}
	}()

	time.Sleep(5 * time.Second)
	close(stopCh) // Отправляем сигнал остановки
	wg.Wait()
}
