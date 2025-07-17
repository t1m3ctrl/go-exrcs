package main

import "fmt"

// Target — это интерфейс, который ожидает клиент
type Target interface {
	Request() string
}

// Тип, несовместимый с интерфейсом
type IncompatibleType struct{}

func (a *IncompatibleType) SpecificRequest() string {
	return "Specific IncompatibleType result"
}

// Адаптер - это обертка над IncompatibleType, реализующая Target
type Adapter struct {
	IncompatibleType *IncompatibleType
}

func (a *Adapter) Request() string {
	return a.IncompatibleType.SpecificRequest()
}

func someCode(t Target) {
	fmt.Println("some code:", t.Request())
}

func main() {
	IncompatibleType := &IncompatibleType{}
	adapter := &Adapter{IncompatibleType: IncompatibleType}
	someCode(adapter)
}
