package main

import (
	"fmt"
)

func main() {
	var s []string
	fmt.Println("Slice:", s, s == nil, len(s) == 0)

	s = make([]string, 3)
	fmt.Println("Initialized Slice:", s, len(s), cap(s))

	s[0] = "a"
	s[1] = "b"
	s[2] = "c"

	fmt.Println("After Assignment:", s)

	s = append(s, "d")

	fmt.Println("After Append:", s)

}
