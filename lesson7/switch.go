package main

import (
	"fmt"
	"time"
)

func main() {
	i := 2

	switch i {
	case 1:
		fmt.Println("One")
	case 2:
		fmt.Println("Two")
	case 3:
		fmt.Println("Three")
	default:
		fmt.Println("Other number")
	}

	switch time.Now().Weekday() {
	case time.Friday:
		fmt.Println("It's Friday, time to relax!")
	case time.Monday:
		fmt.Println("It's Monday, back to work!")
	default:
		fmt.Println("It's just another day.")
	}
}
