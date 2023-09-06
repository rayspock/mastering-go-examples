package main

import "fmt"

func AppendNumber(slice []int) {
	slice = append(slice, 7)
}

func main() {
	slice := make([]int, 3, 6)
	for i := 0; i < len(slice); i++ {
		slice[i] = i
	}
	fmt.Println("before", slice) // Print 0 1 2
	AppendNumber(slice)
	fmt.Println("after", slice) // Print 0 1 2
}
