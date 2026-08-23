package main

import "fmt"

func main() {
	// Nil map
	var nameAge map[string]int
	fmt.Println(nameAge["Alice"]) // 0
	// nameAge["Alice"] = 30         // panic: assignment to entry in nil map

	// Empty map
	nameAge2 := map[string]int{}
	fmt.Println(len(nameAge2))

	// Empty map using make
	nameAge3 := make(map[string]int)
	fmt.Println(len(nameAge3))

	// Empty map with literal variable declaration
	var nameAge4 map[string]int = map[string]int{}
	fmt.Println(len(nameAge4))

	nameAge5 := map[string]int{}
	nameAge5["Bob"] = 25
	nameAge5["Charlie"] = 30
	fmt.Println(len(nameAge5))

	nameAge6 := map[string]int{
		"a": 24,
		"b": 30,
		"c": 28,
	}
	fmt.Println(len(nameAge6))

	// Reading from map
	fmt.Println(nameAge6["a"])   // 24
	fmt.Println(nameAge5["Bob"]) // 25
	fmt.Println(nameAge5["Eve"]) // 0

	// Writing to map
	fmt.Println(nameAge6["c"]) // 28
	nameAge6["c"] = 88
	fmt.Println(nameAge6["c"]) // 88

	// Comma ok idiom
	nameGrade := map[string]int{
		"Alice":   90,
		"Bob":     85,
		"Charlie": 88,
		"David":   0,
	}
	gradeDavid, ok := nameGrade["David"]
	gradeEve, ok2 := nameGrade["Eve"]
	fmt.Println(gradeDavid, ok) // 0 true
	fmt.Println(gradeEve, ok2)  // 0 false

	a := map[string][]int{
		"foo": []int{1, 2, 3},
		"bar": []int{4, 5, 6, 7},
	}
	fmt.Println(a)

}
