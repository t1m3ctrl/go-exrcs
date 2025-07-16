package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"time"
)

func main() {
	args := os.Args[1:]
	N, _ := strconv.Atoi(args[0])

	dataChan := make(chan int)
	timeout := time.After(time.Duration(N) * time.Second)

	go func() {
		for {
			select {
			case <-timeout:
				close(dataChan)
				return
			default:
				num := rand.IntN(100)
				dataChan <- num
				fmt.Printf("Wrote %d\n", num)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	for {
		select {
		case data, ok := <-dataChan:
			if !ok {
				fmt.Println("Channel closed.")
				return
			}
			fmt.Printf("Received: %d\n", data)
		case <-timeout:
			fmt.Println("Timeout, execution is done..")
			return
		}
	}
}
