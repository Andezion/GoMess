package main

import "fmt"

func main() {
	i := 1

	for i <= 5 {
		fmt.Print(i, " ")
		i++
	}

	fmt.Println()

	for _, n := range []int{10, 20, 30, 40, 50} {
		fmt.Print(n, " ")
	}

	fmt.Println()

	for n := range []int{100, 200, 300, 400, 500} {
		fmt.Print(n, " ")
	}
}
