package main

import "fmt"

const x = 10
const y int32 = 15

const (
	m                     = 10
	n               int32 = 15
	applicationName       = "Lesson 7 Constants"
	// All values that can be used as a constant:
	isRunning = true
	character = 'A'
	isTrue    = 1 < 2
)

func main() {
	var a int
	a = x
	fmt.Println(a)

	var b float64
	b = x
	fmt.Println(b)

	var c int
	c = int(y)
	fmt.Println(c)

	var d int32
	d = y
	fmt.Println(d)

	const z = complex(10.2, 100.9)
	const l = imag(z)
	// complex, real, imag, len and cap

}
