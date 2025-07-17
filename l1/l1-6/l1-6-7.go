package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	// Паника и recover
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Восстановлено после:", r)
			}
			wg.Done()
		}()

		fmt.Println("Работает горутина..")
		time.Sleep(500 * time.Millisecond)
		panic("стоп-сигнал")

		// Unreachable
		for {
			fmt.Println("Недостижимый код...")
			time.Sleep(100 * time.Millisecond)
		}
	}()

	time.Sleep(2 * time.Second)
	fmt.Println("Работает main...")

	wg.Wait()
}
