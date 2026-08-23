package main

import (
	"fmt"
	"sync"
)

func main() {
	signal := make(chan any)
	wg := &sync.WaitGroup{}

	for i := range 5 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-signal
			fmt.Printf("Goroutine %d received signal\n", i)
		}(i)
	}
	close(signal)
	wg.Wait()
}
