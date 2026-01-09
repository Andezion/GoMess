package main

import "fmt"

func basic_func(ival int) {
	ival = 0
}

func cool_func(iptr *int) {
	*iptr = 0
}

func main() {
	i := 1
	fmt.Println("initial:", i)

	basic_func(i)
	fmt.Println("zeroval:", i)

	cool_func(&i)
	fmt.Println("zeroptr:", i)

	fmt.Println("pointer:", &i)
}
