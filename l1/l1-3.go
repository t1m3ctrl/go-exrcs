package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"time"
)

func reader(id int, ch <-chan int) {
	for num := range ch {
		fmt.Printf("It's worker %d : %d\n", id, num)
	}
}

func main() {
	args := os.Args[1:]
	ch := make(chan int)

	numWorkers, _ := strconv.Atoi(args[0])

	for i := 0; i < numWorkers; i++ {
		go reader(i, ch)
	}

	for {
		time.Sleep(500 * time.Millisecond)
		num := rand.IntN(100)
		ch <- num
		fmt.Printf("Wrote %d\n", num)
	}
}
