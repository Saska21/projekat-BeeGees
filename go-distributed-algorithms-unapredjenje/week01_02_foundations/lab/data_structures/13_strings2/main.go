package main

import "fmt"

func main() {
	name := "Distributed Algorithms"
	firstLetter := name[0]
	fmt.Println(firstLetter)

	lastLetter := name[len(name)-1]
	fmt.Println(lastLetter)

	// name[0] = 'a' // This will cause a compile-time error: cannot assign to name[0]

	firstLetter2 := name[0]
	var b rune = 'x'
	var c byte = 'y'

	fmt.Println(firstLetter2)
	fmt.Println(b)
	fmt.Println(c)

	runeName := []rune(name)
	fmt.Println(runeName)

	byteName := []byte(name)
	fmt.Println(byteName)
}
