package main

import "fmt"

func main() {
	s := []int{1, 2}

	fmt.Println(s)
	somethingSlice(s)
	fmt.Println(s)

	somethingSlice2(s)
	fmt.Println(len(s))

	myMap := map[string]int{"a": 1, "b": 2}
	fmt.Println(myMap)
	somethingMap(myMap)
	fmt.Println(myMap)
}

// slice internally is implemented as a reference type, which means
// that when you pass a slice to a function,
// you are passing a reference to the underlying array.
func somethingSlice(s []int) {
	s[0] = 10
}

func somethingSlice2(s []int) {
	s = append(s, 1000)
	fmt.Println(len(s))
}

func somethingMap(m map[string]int) {
	delete(m, "a")
}
