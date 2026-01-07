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

	var two_d [2][2]int = [2][2]int{{1, 2}, {3, 4}}
	fmt.Println("Two Dimensional Array:", two_d)

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			fmt.Printf("%d ", two_d[i][j])
		}
		fmt.Println()
	}

	two_d_a := [...][2]int{{5, 6}, {7, 8}}
	fmt.Println("2d: ", two_d_a)

}
