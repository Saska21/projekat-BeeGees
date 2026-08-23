package main

import "fmt"

func main() {
	var intPointer *int
	fmt.Println(intPointer) // Output: <nil>
	//fmt.Println(*intPointer) // This will cause a runtime panic because we are trying to dereference a nil pointer.

	age := 10
	intPointer = &age
	fmt.Println(intPointer)
	fmt.Println(&intPointer)
	// Dereferencing the pointer to get the value
	fmt.Println(*intPointer) // Output: 10

	name := "Alice"
	namePointer := &name
	fmt.Println(namePointer)  // Output: memory address of name
	fmt.Println(*namePointer) // Output: Alice

}
