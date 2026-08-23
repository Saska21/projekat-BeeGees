package main

import (
	"fmt"
	"time"
)

func main() {

	now := time.Now()
	fmt.Println("Current time:", now)

	twoHours := 2 * time.Hour
	futureTime := now.Add(twoHours)
	fmt.Println("Time after 2 hours:", futureTime)

	start := time.Now()
	a := 1
	for i := 0; i < 100000000; i++ {
		a++
	}
	elapse := time.Since(start)
	fmt.Printf("Time taken to execute the loop: %s\n", elapse)
}
