package main

import (
	"fmt"
	"strconv"
)

func main() {
	a := 10
	b := 5
	c := 2.5

	complex1 := complex(2, 3) // 2 + 3i
	complex2 := complex(4, 5) // 4 + 5i

	// sum operator
	sum := a + b
	fmt.Println("Sum:", sum)

	firstName := "Distributed"
	lastName := "Algorithms"
	fullName := firstName + " " + lastName
	fmt.Println("Full Name:", fullName)

	complexSum := complex1 + complex2
	fmt.Println("Complex Sum:", complexSum)

	// Difference operator
	diff := a - b
	fmt.Println("Difference:", diff)

	// Product operator
	product := a * b
	fmt.Println("Product:", product)

	// Division operator
	quotient := a / b
	fmt.Println("Division:", quotient)

	// Remainder operator - can only be used with integers
	remainder := a % b
	fmt.Println("Remainder:", remainder)

	newSum := a + int(c)
	fmt.Println("New Sum:", newSum)

	a += 5
	fmt.Println("Updated a:", a)

	b++
	fmt.Println("Incremented b:", b)

	var d uint8 = 202
	var e uint8 = 141

	// The OR operator
	f := d | e
	fmt.Println("Bitwise OR:", f)
	fmt.Println(strconv.FormatUint(uint64(d), 2)) // binary representation
	fmt.Println(strconv.FormatUint(uint64(e), 2)) // binary representation
	fmt.Println(strconv.FormatUint(uint64(f), 2)) // binary representation

	fmt.Println(11 | 6)

	var g uint8 = 20                                                   // 0001 0100
	fmt.Println(strconv.FormatUint(uint64(g), 2))                      // binary representation
	result := g >> 2                                                   // 0000 0101
	fmt.Println("Right Shift:", strconv.FormatUint(uint64(result), 2)) // binary representation

	g = g << 2                                                   // 0101 0000
	fmt.Println("Left Shift:", strconv.FormatUint(uint64(g), 2)) // binary representation

	h := 1
	h <<= 3
	fmt.Println("Left Shift Assignment:", h)

	h = 5 | (1 << 1)
	fmt.Println("Bitwise OR with Shift:", h)
}
