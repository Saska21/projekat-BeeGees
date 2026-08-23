package main

import "fmt"

func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func main() {
	increment := counter()
	fmt.Println(increment())
	fmt.Println(increment())
	fmt.Println(increment())
	increment2 := counter()
	fmt.Println(increment2())
	fmt.Println(increment2())

	// Anonymous function example
	numbers := []int{1, 2, 3, 4, 5}
	squared := sliceOperation(numbers, func(x int) int {
		return x * x
	})
	fmt.Println("Squared numbers:", squared)
}

func sliceOperation(s []int, op func(int) int) []int {
	result := make([]int, len(s))
	for i, v := range s {
		result[i] = op(v)
	}
	return result
}
