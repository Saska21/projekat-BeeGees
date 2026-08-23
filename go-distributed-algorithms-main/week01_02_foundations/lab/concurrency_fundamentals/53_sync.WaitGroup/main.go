package main

import (
	"fmt"
	"sync"
	"time"
)

func printNumber(wg *sync.WaitGroup, number int) {
	defer wg.Done() // Signal that this goroutine is done
	fmt.Println("Number:", number)
	time.Sleep(3 * time.Second) // Simulate some work

}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go printNumber(&wg, i)
	}

	wg.Wait() // Comment this line to see the difference without WaitGroup

	fmt.Println("All goroutines have finished.")
}
