package main

import "fmt"

var myInteger int = 10

func main() {
	// fmt.Println(myInteger)
	// something()

	// Explicit type assignment
	var myInt int = 10
	var myString string = "Hello, Go!"
	fmt.Println(myInt)
	fmt.Println(myString)

	// Implicit type assignment
	var age = 25
	fmt.Println(age)

	// Multiple variable declaration
	var year, name = 2026, "GoLang"
	fmt.Println(year, name)

	// Short variable declaration
	age2 := 25
	year2 := 2026
	fmt.Println(age2, year2)

	year2, name2 := 2026, "GoLang"
	fmt.Println(year2, name2)

	// Examples:
	/* var apartmentNumber int                             // Zero value to the variable
	var apartmentNumber2 int = 2000                     // When we want to assign a value to variable in declaration
	var apartmentNumber3 = 2000                         // When we want to assign a value to variable in declaration without specifying type
	var apartmentNumber4, streetName4 = 2000, "Main St" // Multiple variable declaration
	apartmentNumber5 := 2000 */
	apartmentNumber6, streetName6 := 2000, "Main St" // Short variable declaration
	fmt.Println(apartmentNumber6, streetName6)
}

func something() {
	fmt.Println("something: ", myInteger)
}
