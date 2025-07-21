package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil // тип = *os.PathError, данные - nil
	return err                  // при возврате преобразовываем в интерфейс error, который хранит
	// все тот же тип и значение nil
}

func main() {
	err := Foo()
	fmt.Println(err)        // data = nil
	fmt.Println(err == nil) // type != nil, type = *os.PathError
	// fmt.Printf("%#v, %#v\n", err, nil)
}
