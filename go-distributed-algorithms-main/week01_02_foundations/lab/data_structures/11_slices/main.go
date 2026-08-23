package main

import "fmt"

func main() {
	// Nil slice
	var a []int
	fmt.Println(a == nil) // true
	// Empty slice
	a = []int{}
	fmt.Println(a == nil) // false

	b := []int{1, 2, 3, 4, 5}
	fmt.Println("Length of b:", len(b))
	fmt.Println("Capacity of b:", cap(b))
	fmt.Println(b)
	b = []int{5, 2: 10, 50}
	fmt.Println(b)

	// Make function to create a slice
	c := make([]int, 5)
	fmt.Println(c)
	fmt.Println("Length of c:", len(c))
	fmt.Println("Capacity of c:", cap(c))

	d := make([]string, 3, 10)
	fmt.Println(d)
	fmt.Println("Length of d:", len(d))
	fmt.Println("Capacity of d:", cap(d))

	// Append function
	e := make([]int, 5)
	e = append(e, 10)
	fmt.Println(e)
	fmt.Println("Length of e:", len(e))
	fmt.Println("Capacity of e:", cap(e))
	e = append(e, 20, 30, 40)
	fmt.Println(e)
	fmt.Println("Length of e:", len(e))
	fmt.Println("Capacity of e:", cap(e))

	someFunction(e)
}

func someFunction(x []int) {}
