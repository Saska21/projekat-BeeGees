package main

import "fmt"

type Rectangle struct {
	width  float64
	height float64
}

func (r *Rectangle) Area() float64 {
	if r == nil {
		return 0
	}
	return r.width * r.height
}

func (r *Rectangle) SetHeight(height float64) {
	r.height = height
}

func updateHeight(r *Rectangle, height float64) {
	r.SetHeight(height)
}

func main() {
	rect := Rectangle{width: 5, height: 5}
	fmt.Println(rect.Area())
	updateHeight(&rect, 10)
	fmt.Println(rect.Area())
}
