package main

import "fmt"

func sum(numbers ...int) int {
	total := 0
	for _, number := range numbers {
		total += number
	}
	return total
}

func main() {
	result := sum(1, 2, 3, 4, 5)
	fmt.Println("Sum:", result)

	result2 := sum(10, 20, 30)
	fmt.Println("Sum:", result2)
}
