package main

import (
	"fmt"
	"sync"
)

var once sync.Once
var initialized bool

func initSomething() {
	fmt.Println("Initializing...")
	initialized = true
}

func doSomething(wg *sync.WaitGroup) {
	once.Do(initSomething)
	once.Do(func() { fmt.Println("This will not run") }) // This will not run because once.Do only executes the first function
	fmt.Println("Doing something...")
	wg.Done()
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(3)

	for _ = range 3 {
		go doSomething(wg)
	}
	wg.Wait()

}
