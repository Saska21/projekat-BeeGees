package main

import (
	"fmt"
	"sync"
)

func main() {
	var counter int
	var lock sync.Mutex

	for range 1_000 {
		go func() {
			lock.Lock()
			counter++
			lock.Unlock()
		}()
	}

	// time.Sleep(10 * time.Second) // Wait for goroutines to finish (not a reliable method in production)

	lock.Lock()
	fmt.Println("Expected counter to be 1_000")
	fmt.Println("Actual counter:", counter)
	lock.Unlock()
}
