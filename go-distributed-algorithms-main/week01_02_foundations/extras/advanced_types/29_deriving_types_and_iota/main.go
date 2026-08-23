package main

import "fmt"

type Age int

// Inheritance of types using type aliasing
type Human struct {
	Name string
	Age  Age
}

func sum(a Age) {
	fmt.Println("Age is:", a)
}

func (h Human) printAge() {
	fmt.Println(h.Age)
}

type Student Human

type Size int

const (
	ExtraSmall Size = iota * 3
	Small
	Medium
	Large
	ExtraLarge
)

func printSize(s Size) {
	switch s {
	case ExtraSmall:
		fmt.Println("Size is Extra Small")
	case Small:
		fmt.Println("Size is Small")
	case Medium:
		fmt.Println("Size is Medium")
	case Large:
		fmt.Println("Size is Large")
	case ExtraLarge:
		fmt.Println("Size is Extra Large")
	default:
		fmt.Println("Unknown Size")
	}
}

func main() {
	var young Age = 10
	var old Age = 60
	fmt.Println(young + old)
	sum(young)

	var s = Student{Name: "Bob", Age: 20}
	// s.printAge() // This will not work because Student does not inherit methods from Human

	// To use the method, we need to convert Student to Human
	h := Human(s)
	h.printAge() // This will work because we converted Student to Human

	fmt.Println("Sizes:", ExtraSmall, Small, Medium, Large, ExtraLarge)
}
