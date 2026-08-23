package main

import (
	"fmt"

	_ "github.com/danilokacanski/init-example/init" // Importing the init package for its side effects
)

func init() {
	fmt.Println("This is the init 1 function.")
}

func init() {
	fmt.Println("This is the init 2 function.")
}

func init() {
	fmt.Println("This is the init 3 function.")
}

func main() {
	fmt.Println("This is the main function.")
}
