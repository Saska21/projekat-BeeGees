package main

import "fmt"

func main() {
	var smallFloat float32
	fmt.Println(smallFloat)
	smallFloat = 23.0234324
	fmt.Println(smallFloat)

	var bigFloat float64
	fmt.Println(bigFloat)
	bigFloat = 234234234234234.234234234234
	fmt.Println(bigFloat)

	var myComplex complex128
	myComplex = complex(bigFloat, bigFloat)
	fmt.Println(myComplex)

	var myRealPart, myComplexPart float64
	myRealPart = real(myComplex)
	myComplexPart = imag(myComplex)
	fmt.Println(myRealPart)
	fmt.Println(myComplexPart)

}
