package main

import "fmt"

type Human struct {
	age  int
	name string
}

func (h Human) say() {
	fmt.Println("I'm Human")
}

func (h Human) cheers() {
	fmt.Println("I'm Human and i'm cheering")
}

type Action struct {
	Human
	action string
}

func (a Action) say() {
	fmt.Println("I'm action")
}

func main() {
	a := Action{}

	a.say()
	a.cheers()
	a.Human.say()
}
