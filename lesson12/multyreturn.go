package main

import "fmt"

func return_nuumbers() (int, int) {
	return 5, 10
}

func main() {
	num1, num2 := return_nuumbers()
	fmt.Println("First number:", num1)
	fmt.Println("Second number:", num2)
}
