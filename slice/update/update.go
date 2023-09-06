package main

import "fmt"

func UpdateNumber(slice []int) {
	slice[1] = 7
}

func main() {
	slice := make([]int, 3, 6)
	for i := 0; i < len(slice); i++ {
		slice[i] = i
	}
	fmt.Println("before", slice) // Print 0 1 2
	UpdateNumber(slice)
	fmt.Println("after", slice) // Print 0 7 2
}
