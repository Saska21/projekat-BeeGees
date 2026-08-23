package main

import "fmt"

type Circle struct {
	radius float64
}

type Rectangle struct {
	width  float64
	height float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.radius * c.radius
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}

func (r *Rectangle) SetHeight(height float64) { // Remove the pointer to see effect
	r.height = height
}

func main() {
	rect := Rectangle{width: 5, height: 10}
	fmt.Println("Area of rectangle:", rect.Area())
	rect.SetHeight(20)
	fmt.Println("New area of rectangle:", rect.Area())

}
