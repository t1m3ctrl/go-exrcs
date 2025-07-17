package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Остановка контекстом")
				return
			default:
				fmt.Println("Работает...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	time.Sleep(5 * time.Second)
	cancel() // Отмена контекста
	wg.Wait()
}
