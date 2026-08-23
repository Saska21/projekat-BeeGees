package main

import "fmt"

func main() {
	age := 10
	if age >= 18 {
		fmt.Println("Adult")
	} else if age >= 13 {
		fmt.Println("Teenager")
	} else {
		fmt.Println("Child")
	}

	if even := isEven(age); even == true {
		fmt.Println("Even age")
	}
	// fmt.Println(even) // This will cause a compile-time error: undefined: even
}
func isEven(n int) bool {
	return n&1 == 0
}
