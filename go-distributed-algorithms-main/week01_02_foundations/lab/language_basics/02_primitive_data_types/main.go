package main

import "fmt"

func main() {
	var myString string
	// String zero value is ""
	fmt.Println(myString)
	myString = "Welcome to Go!"
	fmt.Println(myString)

	var myInt int
	// Zero value of number types is 0
	fmt.Println(myInt)
	myInt = 42
	fmt.Println(myInt)

	var myBool bool
	// Zero value of boolean type is false
	fmt.Println(myBool)
	myBool = true
	fmt.Println(myBool)
}
