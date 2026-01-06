package main

import "fmt"

func main() {
	var nigga bool = true

	if nigga == true {
		fmt.Println("u r true nigga")
	} else {
		fmt.Println("u r bitch ass")
	}

	if num := 9; num < 0 {
		fmt.Println(num, "is negative")
	} else if num < 10 {
		fmt.Println(num, "has 1 digit")
	} else {
		fmt.Println(num, "has multiple digits")
	}
}
