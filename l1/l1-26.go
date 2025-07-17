package main

import (
	"fmt"
	"strings"
)

func isUnique(s string) bool {
	s = strings.ToLower(s)
	freq := map[rune]int{}
	for _, v := range s {
		freq[v]++
		if freq[v] > 1 {
			return false
		}
	}
	return true
}

func main() {
	str1 := "abCdefAaf"
	str2 := "abcd"
	str3 := "aabcd"

	fmt.Println("str1 is unique? ->", isUnique(str1))
	fmt.Println("str2 is unique? ->", isUnique(str2))
	fmt.Println("str3 is unique? ->", isUnique(str3))
}
