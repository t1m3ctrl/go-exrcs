package main

import "fmt"

func reverse(s string) string {
	runes := []rune(s)

	// свап конца с началом до середины
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func main() {
	fmt.Println(reverse("главрыба"))
	fmt.Println(reverse("hello 🌍🚀"))
}
