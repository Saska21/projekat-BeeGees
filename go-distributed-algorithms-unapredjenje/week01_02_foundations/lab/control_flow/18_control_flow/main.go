package main

import "fmt"

func main() {
	a := 100

	for {
		fmt.Println(a)
		a--
		if a == 90 {
			break
		}
	}

	for i := 0; i < 5; i++ {
		if i == 2 {
			continue
		}
		fmt.Println(i)
	}

	// Labeled break
outerLoop:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				break outerLoop
			}
			fmt.Printf("i: %d, j: %d\n", i, j)
		}
	}

	b := 10

	if b == 10 {
		goto end
	}
	// c := 20
	fmt.Println("This will be skipped")
end:
	fmt.Println("End of main function")

}
