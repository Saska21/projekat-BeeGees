package main

import "fmt"

func main() {
	//	uint8
	var smallPositiveValue uint8
	smallPositiveValue = 10
	fmt.Println(smallPositiveValue)
	// int8
	var smallPositiveNegativeValue int8
	smallPositiveNegativeValue = -10
	fmt.Println(smallPositiveNegativeValue)
	// uint16, uint32, uint64
	// int16, int32, int64
	var myInt int = 23000043
	fmt.Println(myInt)

	myInt = int(smallPositiveNegativeValue)
	myInt = int(smallPositiveValue)
	// to typecase a value int()

	// byte is mainly used ot represent raw data
	// and also distinguish between uint8 and byte
	// since byte is an alias for uint8
	var myByte byte = 'A'
	fmt.Println(myByte)

	// rune is an alias for int32
	var myRune rune = 'B'
	fmt.Println(myRune)
	myRune = '😄' // can't add space here
	fmt.Println(myRune)
}
