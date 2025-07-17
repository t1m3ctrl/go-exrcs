package main

import (
	"fmt"
)

func reverseRunes(r []rune, left, right int) {
	for left < right {
		r[left], r[right] = r[right], r[left]
		left++
		right--
	}
}

func reverseWords(s string) string {
	r := []rune(s)
	n := len(r)

	reverseRunes(r, 0, n-1)

	start := 0
	for i := 0; i <= n; i++ {
		if i == n || r[i] == ' ' {
			reverseRunes(r, start, i-1)
			start = i + 1
		}
	}

	return string(r)
}

func main() {
	input := "snow dog sun"
	output := reverseWords(input)
	fmt.Println(output)

	input2 := "снег собака солнце"
	output2 := reverseWords(input2)
	fmt.Println(output2)
}
