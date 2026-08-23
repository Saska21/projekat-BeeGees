package main

import "fmt"

type multiplier func(int) int

func multiplyBy(multiplier int) multiplier {
	return func(i int) int {
		return i * multiplier
	}
}

type operation func(int, int) int

func arithemticOperation(op string) operation {
	switch op {
	case "add":
		return func(a, b int) int {
			return a + b
		}
	case "subtract":
		return func(a, b int) int {
			return a - b
		}
	default:
		return nil
	}
}

func main() {
	multiplyByTwo := multiplyBy(2)
	multiplyByThree := multiplyBy(3)

	result1 := multiplyByTwo(5)
	result2 := multiplyByThree(5)

	fmt.Println("5 multiplied by 2 is:", result1)
	fmt.Println("5 multiplied by 3 is:", result2)

	var perform operation
	perform = arithemticOperation("add")
	result := perform(1, 2)
	fmt.Println(result)
}
