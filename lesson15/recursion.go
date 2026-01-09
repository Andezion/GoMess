package main

import "fmt"

func factorial(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}

func main() {
	fmt.Println(factorial(5))

	var test func(n int) int

	test = func(n int) int {
		if n == 0 {
			return 1
		}
		return n * test(n-1)
	}

	fmt.Println(test(7))
}
