package main

import (
	"fmt"
	"math"

	internal "github.com/danilokacanski/da/package-examples/math" // Importing the math package from the internal directory with alias
)

func main() {
	sum := internal.Add(5, 3)
	fmt.Println("Sum:", sum)

	PI := internal.PI
	fmt.Println("PI:", PI)

	// The following line would cause an error because subtract is unexported
	// difference := internal.subtract(5, 3)
	// fmt.Println("Difference:", difference)

	math.Abs(20.0)
}
