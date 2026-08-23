package main

import "fmt"

func main() {
	integerChannel := make(chan int)
	// writeOnlyChannel(integerChannel)

	// Part 1
	// go func() {
	// 	integerChannel <- 42
	// }()

	// readValue, ok := <-integerChannel
	// fmt.Println(readValue, ok)

	// close(integerChannel)

	// readValue, ok = <-integerChannel
	// fmt.Println(readValue, ok)

	// Part 2
	go func() {
		for i := range 5 {
			integerChannel <- i
		}
		close(integerChannel)
	}()

	for j := range integerChannel {
		fmt.Println(j)
	}

}

func readOnlyChannel(ch <-chan int)  {}
func writeOnlyChannel(ch chan<- int) {}
