package main

func main() {
	ch := make(chan int)
	go func() {
		// defer close(ch)
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()
	for n := range ch { // можно читать после закрытия
		println(n)
	}
}
