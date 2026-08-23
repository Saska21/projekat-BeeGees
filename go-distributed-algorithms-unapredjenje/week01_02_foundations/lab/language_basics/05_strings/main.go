package main

import "fmt"

func main() {
	var myString string
	fmt.Println(myString)

	myString = "Hello, Go!"
	fmt.Println(myString)

	myString = "Hello, \nGo!"
	fmt.Println(myString)

	var firstName, lastName string
	firstName = "Distributed"
	lastName = "Algorithms"

	var fullName string
	fullName = firstName + " " + lastName
	fmt.Println(fullName)
}
