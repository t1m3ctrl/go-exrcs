package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v, ok := <-a:
				if ok {
					c <- v
				} else {
					a = nil
				}
			case v, ok := <-b:
				if ok {
					c <- v
				} else {
					b = nil
				}
			}
			if a == nil && b == nil {
				close(c)
				return
			}
		}
	}()
	return c
}

func main() {
	rand.Seed(time.Now().Unix())
	a := asChan(1, 3, 5, 7) // в каналы a и b с периодичностью пишутся данные
	b := asChan(2, 4, 6, 8) // все данные в канал b МОГУТ БЫТЬ записаны быстрее чем второе значение в канал a,
	// поэтому вывод может быть таким: 1, 2, 3, 5, 7, 4, 6, 8
	c := merge(a, b) // эти данные читаются из каналов и объединяются в один, до тех пор пока данные не закончатся
	// или каналы не закроются
	for v := range c { // читаются данные из канала c
		fmt.Print(v)
	}
}
