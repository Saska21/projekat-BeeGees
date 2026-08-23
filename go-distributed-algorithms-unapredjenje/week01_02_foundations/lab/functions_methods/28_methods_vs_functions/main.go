package main

import "fmt"

func main() {
	calc := Calculator{}
	fmt.Println(calc.Add(2, 3)) // Using the method on the struct

	fmt.Println(AddFunction(2, 3)) // Using the function directly

	fmt.Println(ArithemeticOperation(AddFunction, 2, 3))

	printer := &Printer{}

	printFunction := printer.Print // Assigning the method to a variable (function value)

	printFunction("Hello, World!") // Using the returned function to print

	f1 := Person.GetDetails

	p := Person{Name: "Alice", Age: 30}

	fmt.Println(f1(p))

}
