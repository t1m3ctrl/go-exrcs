package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func reader(ctx context.Context, id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case num, ok := <-ch:
			if !ok {
				fmt.Printf("Worker %d: channel closed, exiting\n", id)
				return
			}
			fmt.Printf("It's worker %d : %d\n", id, num)
		case <-ctx.Done():
			fmt.Printf("Worker %d: cancellation signal received, exiting\n", id)
			return
		}
	}
}

func main() {
	args := os.Args[1:]
	ch := make(chan int)

	numWorkers, _ := strconv.Atoi(args[0])

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go reader(ctx, i, ch, &wg)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

loop:
	for {
		select {
		case <-sigCh:
			fmt.Println("\nSIGINT received, shutting down gracefully...")
			cancel()
			close(ch)
			break loop
		default:
			time.Sleep(500 * time.Millisecond)
			num := rand.IntN(100)
			ch <- num
			fmt.Printf("Wrote %d\n", num)
		}
	}

	wg.Wait()
	fmt.Println("All workers exited, program terminated")
}
