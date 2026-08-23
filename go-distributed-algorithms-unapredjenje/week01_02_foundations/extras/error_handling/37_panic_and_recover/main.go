package main

import "fmt"

func function1() {
	defer func() {
		fmt.Println("Function 1: Deferred function called")
	}()
	function2()
}

func function2() {
	defer func() {
		fmt.Println("Function 2: Deferred function called")
	}()
	panic("Function 2: Something went wrong!")
}

func main() {
	// defer func() {
	// 	fmt.Println("Main: Deferred function called")
	// }()
	// function1()

	fmt.Println("Starting of the program")

	panicExample()

	fmt.Println("End of the program")
}
