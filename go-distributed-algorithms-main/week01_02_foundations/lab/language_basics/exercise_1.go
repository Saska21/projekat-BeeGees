package main

import "fmt"

func main() {
	// Variables: string, int, bool, float64, rune
	var course string = "Go Basics"
	var students int = 24
	var isLab bool = true
	var score float64 = 92.5
	var symbol rune = 'G'

	// Print values and types
	fmt.Printf("course=%v, type=%T\n", course, course)
	fmt.Printf("students=%v, type=%T\n", students, students)
	fmt.Printf("isLab=%v, type=%T\n", isLab, isLab)
	fmt.Printf("score=%v, type=%T\n", score, score)
	fmt.Printf("symbol=%v, type=%T\n", symbol, symbol)

	// Constants: one untyped and one typed
	const maxPoints = 100       // untyped constant
	const passMark float64 = 51 // typed constant
	fmt.Printf("maxPoints=%v, passMark=%v\n", maxPoints, passMark)

	// Type conversion: int -> float64
	studentsAsFloat := float64(students)
	fmt.Printf("studentsAsFloat=%v, type=%T\n", studentsAsFloat, studentsAsFloat)

	// String concatenation
	message := "Welcome to " + course
	fmt.Println(message)

	// Bitwise operation: left shift
	var mask uint8 = 1
	shifted := mask << 3
	fmt.Printf("mask=%08b, shifted=%08b\n", mask, shifted)

	// Bonus runtime panic example (commented out)
	// var p *int
	// fmt.Println(*p) // panic: invalid memory address or nil pointer dereference
}
