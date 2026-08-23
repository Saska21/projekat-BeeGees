package main

import (
	"fmt"
	"sync"
	"time"
)

var cond = sync.NewCond(&sync.Mutex{})
var integers = make([]int, 0, 10)

func remove(delay time.Duration) {
	time.Sleep(delay)
	fmt.Println("Len before removing:", len(integers))
	integers = integers[1:]
	fmt.Println("Len after removing:", len(integers))
	cond.Broadcast()
}

func add() {
	for i := 0; i < 5; i++ {
		cond.L.Lock()
		for len(integers) == 2 {
			cond.Wait()
			fmt.Println("Len before adding:", len(integers))
		}
		integers = append(integers, i)
		fmt.Println("Len after adding:", len(integers))
		go remove(1 * time.Second)
		cond.L.Unlock()
	}
}

func main() {
	go add()
	time.Sleep(5 * time.Second)
}
