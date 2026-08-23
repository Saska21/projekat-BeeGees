package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for i := range 3 {
			ch1 <- i
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for i := range 3 {
			ch2 <- i
			time.Sleep(1 * time.Second)
		}
	}()
	for range 6 {
		select {
		case value := <-ch1:
			fmt.Println("Received", value, "from ch1")
		case value := <-ch2:
			fmt.Println("Received", value, "from ch2")
			// case <-time.After(1 * time.Second):
			// 	fmt.Println("Timeout: no data received in the last second")
			// 	return
		}
	}
}
