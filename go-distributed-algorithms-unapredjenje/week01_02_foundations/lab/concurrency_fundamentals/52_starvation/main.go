package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	mutex   sync.Mutex
	counter int
)

func highPriority() {
	mutex.Lock()
	time.Sleep(100 * time.Millisecond)
	defer mutex.Unlock()
	counter += 10
}

func lowPriority() {
	mutex.Lock()
	defer mutex.Unlock()
	counter++
}

func main() {
	for i := 0; i < 1_000; i++ {
		go highPriority()
		go lowPriority()
	}
	time.Sleep(10 * time.Second) // Wait for goroutines to finish (not a reliable method in production)
	fmt.Println("Counter:", counter)
}
