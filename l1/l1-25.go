package main

import (
	"fmt"
	"time"
)

func sleep(d time.Duration) {
	timer := time.NewTimer(d)
	<-timer.C
}

func main() {
	fmt.Println("Hello ...")
	sleep(2 * time.Second)
	fmt.Println("... World!")
}
