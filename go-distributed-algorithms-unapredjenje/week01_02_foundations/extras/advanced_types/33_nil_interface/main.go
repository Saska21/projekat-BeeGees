package main

import "fmt"

func main() {
	var a *int
	var i interface{}

	fmt.Println(a == nil) // true
	fmt.Println(i == nil) // true

	i = a
	fmt.Println(i == nil) // false, because i holds a value of type *int, even though that value is nil
}
