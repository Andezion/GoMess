package main

import "fmt"

func test() func() int {
	x := 0
	return func() int {
		x++
		return x
	}
}

func main() {
	res := test()

	fmt.Println(res()) // 1
	fmt.Println(res()) // 2
	fmt.Println(res()) // 3

	res1 := test()
	fmt.Println(res1()) // 1
}
