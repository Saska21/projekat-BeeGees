package main

import "fmt"

type Shape interface {
	Area() float64
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Name() string {
	return "Rectangle"
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.Radius * c.Radius
}

func printArea(s Shape) {
	fmt.Println("Area:", s.Area())
}

type Square struct {
	Side float64
}

func (s Square) Area() float64 {
	return s.Side * s.Side
}

func (s Square) PrintArea() {
	printArea(s)
}

type Object interface {
	Name() string
	Shape
}

func printObject(o Object) {
	fmt.Printf("Area: %f, Name: %s", o.Area(), o.Name())
}

func main() {
	// Go does not have concept of interfaces, but we can achieve similar results using type embedding and method promotion
	rectangle := Rectangle{Width: 5, Height: 10}
	circle := Circle{Radius: 7}
	square := Square{Side: 4}

	shapes := []Shape{rectangle, circle}

	for _, shape := range shapes {
		printArea(shape)
	}

	printArea(square) // Square also implements the Shape interface, so we can use it here
	square.PrintArea()

	// Embedding with interfaces
	printObject(rectangle)
}
