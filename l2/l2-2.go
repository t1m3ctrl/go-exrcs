package main

import "fmt"

func test() (x int) {
	defer func() {
		x++ // увеличивает именнованный параметр
	}()
	x = 1
	return x // x = 1, defer x++, x = 2
}

func anotherTest() int {
	var x int
	defer func() {
		x++ // не оказывает влияния на возвращаемое значение
	}()
	x = 1
	return x // На данном этапе уже зафиксировано, что будет возвращено x = 1
}

func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
