package main

import (
	"fmt"
)

func main() {
	var arr [5]int
	fmt.Println("Array:", arr)

	b := [5]int{1, 2, 3, 4, 5}
	fmt.Println("Initialized Array:", b)

	c := [...]int{10, 20, 30}
	fmt.Println("Inferred Size Array:", c)
	fmt.Println("Size:", len(c))
}
