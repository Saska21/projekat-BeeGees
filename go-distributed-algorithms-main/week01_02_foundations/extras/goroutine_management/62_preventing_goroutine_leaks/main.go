package main

import (
	"fmt"
	"sync"
	"time"
)

func someWork(done <-chan struct{}, wg *sync.WaitGroup) <-chan int {
	ch := make(chan int)

	go func() {
		defer wg.Done()
		select {
		case ch <- 1:
		case <-done:
			return
		}
		fmt.Println("Work done")
	}()

	return ch
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	done := make(chan struct{})
	someWork(done, wg)
	time.Sleep(1 * time.Second)
	close(done)
	wg.Wait()
	fmt.Println("All work completed")
}
