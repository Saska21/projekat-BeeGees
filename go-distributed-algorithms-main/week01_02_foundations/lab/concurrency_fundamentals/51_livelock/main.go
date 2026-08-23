package main

import (
	"fmt"
	"sync"
)

var (
	mutexA sync.Mutex
	mutexB sync.Mutex
)

func routineA(name string) {
	mutexA.Lock()
	defer mutexA.Unlock()
	fmt.Println(name, "has mutexA")
	fmt.Println(name, "trying to acquire mutexB...")
	if !mutexB.TryLock() {
		fmt.Println(name, "could not acquire mutexB, releasing mutexA and retrying...")
		// Simulate some work before retrying
		for i := 0; i < 1_000; i++ {
			// Simulating work
		}
		routineA(name) // Retry acquiring locks
	} else {
		defer mutexB.Unlock()
		fmt.Println(name, "acquired mutexB")
		fmt.Println(name, "doing some work...")
	}
}

func routineB(name string) {
	mutexB.Lock()
	defer mutexB.Unlock()
	fmt.Println(name, "has mutexB")
	fmt.Println(name, "trying to acquire mutexA...")
	if !mutexA.TryLock() {
		fmt.Println(name, "could not acquire mutexA, releasing mutexB and retrying...")
		// Simulate some work before retrying
		for i := 0; i < 1_000; i++ {
			// Simulating work
		}
		routineB(name) // Retry acquiring locks
	}
	defer mutexA.Unlock()
	fmt.Println(name, "acquired mutexA")
	fmt.Println(name, "doing some work...")
}

func main() {
	go routineA("Goroutine A")
	go routineB("Goroutine B")

	// Wait for goroutines to finish (not a reliable method in production)
	fmt.Println("Waiting for goroutines to finish...")
	select {}
}
