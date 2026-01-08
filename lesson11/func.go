package main

import "fmt"

func add_two_numbers(a int, b int) int {
	return a + b
}

func add_three_numbers(a int, b int, c int) int {
	return a + b + c
}

func main() {
	sum1 := add_two_numbers(3, 5)
	sum2 := add_three_numbers(2, 4, 6)
	fmt.Println("Sum of two numbers:", sum1)
	fmt.Println("Sum of three numbers:", sum2)
}
