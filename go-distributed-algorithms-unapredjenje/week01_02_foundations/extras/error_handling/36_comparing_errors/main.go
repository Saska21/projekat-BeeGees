package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

var ErrFirstError = CustomError{Message: "First error"}

type CustomError struct {
	Message string
	Status  int
}

func (c CustomError) Error() string {
	return c.Message
}

func someFunction() error {
	return fmt.Errorf("someFunction failed: %w", ErrFirstError)
}

func someOtherFunction() error {
	return CustomError{Message: "someOtherFunction failed", Status: 100}
}

func main() {
	// err := someFunction()
	// fmt.Println(err)
	// if err != nil {
	// 	if errors.Is(err, ErrFirstError) {
	// 		fmt.Println("Sentinel error found")
	// 	}
	// }

	// if _, err := os.Open("non-existing"); err != nil {
	// 	if errors.Is(err, fs.ErrNotExist) {
	// 		fmt.Println("File does not exist")
	// 	} else {
	// 		fmt.Println("An error occurred:", err)
	// 	}
	// }

	// err := someOtherFunction()

	// var customErr CustomError
	// if errors.As(err, &customErr) {
	// 	fmt.Printf("Custom error found: %s with status %d\n", customErr.Message, customErr.Status)
	// } else {
	// 	fmt.Println("An error occurred:", err)
	// }

	if _, err := os.Open("non-existing"); err != nil {
		var pathError *fs.PathError
		if errors.As(err, &pathError) {
			fmt.Printf("Path error: %s\n", pathError)
		}
	}
}
