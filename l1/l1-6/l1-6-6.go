package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	dataCh := make(chan int)

	// та же остановка через канал, только он использовался для передачи
	// данных
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case data, ok := <-dataCh:
				if !ok {
					fmt.Println("Канал закрыт")
					return
				}
				fmt.Println("Обработка данных:", data)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	for i := 0; i < 10; i++ {
		dataCh <- i
	}

	close(dataCh)
	wg.Wait()
}
