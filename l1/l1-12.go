package main

import "fmt"

func main() {
	animals := []string{"cat", "cat", "dog", "cat", "tree"}

	set := map[string]bool{}

	for _, v := range animals {
		set[v] = true
	}

	for k := range set {
		fmt.Printf("%s, ", k)
	}
}
