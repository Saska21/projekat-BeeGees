package main

import (
	"fmt"
)

func main() {
	message := "Hello, "
	greetingFn := func(name string) {
		fmt.Println(message + name)
	}
	// Deferred function calls are executed in LIFO order after the surrounding function returns.
	defer greetingFn("Alice")
	defer greetingFn("Bob")
	defer greetingFn("Charlie")
	fmt.Println("Test")

	// os.Exit(0) // This will prevent deferred functions from executing, demonstrating that they are not called if the program exits before reaching the end of the main function.
	message = "Hi, "
}
