package main

import (
	"fmt"
)

func main() {
	var ch1, ch2 chan any
	select {
	case value := <-ch1:
		fmt.Printf("Received %d from ch1\n", value)
	case ch2 <- struct{}{}:
		fmt.Println("Written to ch2")
	}
}
