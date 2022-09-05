package main

import (
	"bufio"
	"fmt"
	"os"
	"sync/atomic"
)

var ops int32

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Need at least one phrase and file parameters!")
		return
	}
	phrase := os.Args[1]
	A := make(chan int)
	jobNumbers := len(os.Args) - 2
	atomic.AddInt32(&ops, int32(jobNumbers)) // to keep track of total numbers of goroutines
	for i := 2; i < len(os.Args); i++ {
		fn := os.Args[i]
		// Reads the file and finds the number of occurrences of a given phrase
		go scanFile(fn, phrase, A)
	}
	// Calculate total numbers of occurrences
	total := sum(A)
	fmt.Println("total occurrences:", total)
}

func sum(in <-chan int) int {
	total := 0
	for x := range in {
		total += x
	}
	return total
}

func scanFile(fn, phrase string, out chan<- int) {
	defer func() {
		atomic.AddInt32(&ops, -1)
		if atomic.LoadInt32(&ops) == 0 {
			close(out)
			return
		}
	}()
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	count := 0
	for scanner.Scan() {
		if phrase == scanner.Text() {
			count++
		}
	}
	fmt.Println("file:", f.Name(), ", occurrences:", count)
	// send to the channel
	out <- count
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}
}
