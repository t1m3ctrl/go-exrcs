package main

import (
	"fmt"
)

func main() {
	var s = []string{"1", "2", "3"}
	modifySlice(s)
	fmt.Println(s)
}

func modifySlice(i []string) {
	// fmt.Printf("Начало: len=%d, cap=%d, ptr=%p, data %v\n", len(i), cap(i), &i[0], i)
	i[0] = "3" // {3, 2, 3}
	// fmt.Printf("Изменили i[0]: len=%d, cap=%d, ptr=%p, data %v\n", len(i), cap(i), &i[0], i)
	i = append(i, "4") // new i (local) {3, 2, 3, 4}
	// fmt.Printf("Добавили \"4\": len=%d, cap=%d, ptr=%p, data %v\n", len(i), cap(i), &i[0], i)
	i[1] = "5" // with new i (local) {3, 5, 3, 4}
	// fmt.Printf("Изменили i[1]: len=%d, cap=%d, ptr=%p, data %v\n", len(i), cap(i), &i[0], i)
	i = append(i, "6") // change the new i (local) {3, 5, 3, 4, 6}
	// fmt.Printf("Добавили \"6\": len=%d, cap=%d, ptr=%p, data %v\n", len(i), cap(i), &i[0], i)
}
