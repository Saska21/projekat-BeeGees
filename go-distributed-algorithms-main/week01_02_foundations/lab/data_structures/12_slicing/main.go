package main

import "fmt"

func main() {
	a := []int{1, 2, 3, 4, 5, 6}
	fmt.Println(a)

	e := a[2:6:6]

	b := a[2:5]
	fmt.Println(b)

	c := a[:3]
	d := a[3:]

	fmt.Println(c)
	fmt.Println(d)

	// Ripple effect
	c[0] = 100
	fmt.Println(a)
	fmt.Println(c)

	b = append(b, 200)
	fmt.Println(a)
	fmt.Println(b)

	fmt.Println(e)

	// No ripple effect
	f := [6]int{10, 20, 30, 40, 50, 60}
	g := make([]int, 6)
	x := copy(g, f[:])
	fmt.Println(x)
	fmt.Println(g)

}
