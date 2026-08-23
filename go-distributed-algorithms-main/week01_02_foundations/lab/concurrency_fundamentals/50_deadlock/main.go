package main

import (
	"fmt"
	"sync"
)

var (
	mutex1 sync.Mutex
	mutex2 sync.Mutex
)

func resource1(name string) {
	mutex1.Lock()
	fmt.Println(name, "acquired resource 1")
	resource2(name)
	mutex1.Unlock()
}

func resource2(name string) {
	mutex2.Lock()
	fmt.Println(name, "acquired resource 2")
	resource1(name)
	mutex2.Unlock()
}

func main() {
	go resource1("Goroutine 1")
	go resource2("Goroutine 2")

	// Wait for goroutines to finish (not a reliable method in production)
	fmt.Println("Waiting for goroutines to finish...")
	select {}

}
