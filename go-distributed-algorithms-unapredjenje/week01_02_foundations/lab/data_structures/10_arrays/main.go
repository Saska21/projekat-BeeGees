package main

import "fmt"

func main() {
	var a [5]int
	fmt.Println(a)

	b := [5]int{1, 2, 3, 4, 5}
	fmt.Println(b)

	// Sparse array declaration; 0th and 3rd indices are initialized
	c := [5]int{0: 10, 3: 40}
	fmt.Println(c)

	// Implicit lenght declaration
	d := [...]int{3, 4, 1, 2}
	fmt.Println(d)
	fmt.Println("Length of d:", len(d))
	e := d[2]
	fmt.Println(e)

	f := [2][2]int{{1, 2}, {3, 4}}
	fmt.Println(f[0][1])

}
