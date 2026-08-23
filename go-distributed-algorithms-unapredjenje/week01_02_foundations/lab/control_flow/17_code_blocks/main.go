package main

import "fmt"

var (
	// Package-level variable
	a = 10
)

func main() {
	fmt.Println(a)
	{
		// Code block scope
		a := 15
		fmt.Println(a)
		something()
	}
	fmt.Println(a)
	something()
}
