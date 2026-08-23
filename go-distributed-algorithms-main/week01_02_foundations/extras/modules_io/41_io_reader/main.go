package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {

	// f, err := os.Open("letters.txt")
	// if err != nil {
	// 	panic(err)
	// }

	// n, err := countAlphabets(f)
	// if err != nil {
	// 	panic(err)
	// }

	// println("Number of alphabets in the file:", n)

	r := strings.NewReader("Hello, Worldddd$$$$!")
	n, err := countAlphabets(r)
	if err != nil {
		panic(err)
	}
	fmt.Println("Number of alphabets in the string:", n)
}

func countAlphabets(r io.Reader) (int, error) {
	count := 0
	buffer := make([]byte, 1024)

	for {
		n, err := r.Read(buffer)
		for _, l := range buffer[:n] {
			if (l >= 'a' && l <= 'z') || (l >= 'A' && l <= 'Z') {
				count++
			}
		}
		if err == io.EOF {
			return count, nil
		}
		if err != nil {
			return 0, err
		}
	}
}
