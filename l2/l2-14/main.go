package main

import (
	"fmt"
	"sync"
	"time"
)

// or объединяет один или более каналов done (каналов сигнала завершения) в один,
// который закрывается, как только закроется любой из исходных каналов
func or(channels ...<-chan interface{}) <-chan interface{} {
	if len(channels) == 0 {
		return nil
	}

	if len(channels) == 1 {
		return channels[0]
	}

	result := make(chan interface{})
	var done sync.Once

	for _, ch := range channels {
		go func(c <-chan interface{}) {
			select {
			case <-c:
				done.Do(func() {
					close(result)
				})
			case <-result:
			}
		}(ch)
	}

	return result
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}
