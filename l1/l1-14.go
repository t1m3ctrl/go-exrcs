package main

import (
	"fmt"
	"reflect"
)

func guessType(x interface{}) string {
	switch x.(type) {
	case int:
		return "int"
	case string:
		return "string"
	case bool:
		return "bool"
	case chan interface{}:
		return "channel"
	default:
		t := fmt.Sprintf("%T", x)
		if len(t) > 4 && t[:4] == "chan" {
			return "channel"
		}
		return "Don't know about this type"
	}
}

func guessType2(x interface{}) string {
	switch reflect.TypeOf(x).Kind() {
	case reflect.Int:
		return "int"
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "bool"
	case reflect.Chan:
		return "channel"
	default:
		return "Don't know about this type"
	}
}

func main() {
	var a int
	var name string
	var condition bool
	var channel chan int
	var channelString chan string

	var unknownType float32

	fmt.Println(guessType(a))
	fmt.Println(guessType(name))
	fmt.Println(guessType(condition))
	fmt.Println(guessType(channel))
	fmt.Println(guessType(channelString))
	fmt.Println(guessType(unknownType))
	fmt.Println()
	fmt.Println(guessType2(a))
	fmt.Println(guessType2(name))
	fmt.Println(guessType2(condition))
	fmt.Println(guessType2(channel))
	fmt.Println(guessType2(channelString))
	fmt.Println(guessType2(unknownType))
}
