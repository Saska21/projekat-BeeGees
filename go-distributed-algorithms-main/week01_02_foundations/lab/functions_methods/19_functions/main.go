package main

import "fmt"

func main() {
	sayHi()
	sayHiToSomeone("Distributed Algorithms")
	fn := fullName("Distributed", "Algorithms")
	fmt.Println(fn)
	fn2, length := fullNameWithLength("Go", "Lang")
	fmt.Printf("Full Name: %s, Length: %d\n", fn2, length)

	// Calling varadic function
	numbers := []int{1, 2, 3, 4, 5, 6}
	s := sum(numbers...)
	fmt.Printf("Sum: %d\n", s)

	// Function signatures
	result := apply(double, 5)
	fmt.Println("Result of apply:", result)
}

func sayHi() {
	fmt.Println("Hi")
}

func sayHiToSomeone(name string) {
	fmt.Printf("Hi, %s!\n", name)
}

func fullName(firstName, lastName string) string {
	return fmt.Sprintf("%s %s", firstName, lastName)
}

func fullNameWithLength(firstName, lastName string) (string, int) {
	fullName := fullName(firstName, lastName)
	return fullName, len(fullName)
}

// Variadic function example
func sum(numbers ...int) int {
	total := 0
	for _, number := range numbers {
		total += number
	}
	return total
}

func double(x int) int {
	return x * 2
}

func apply(f func(int) int, value int) int {
	return f(value)
}
