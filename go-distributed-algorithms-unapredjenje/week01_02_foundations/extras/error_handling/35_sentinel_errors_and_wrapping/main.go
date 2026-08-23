package main

import (
	"errors"
	"fmt"
)

func firstFunction() error {
	return fmt.Errorf("this is an error from firstFunction")
}

func secondFunction() error {
	firstErr := firstFunction()
	if firstErr != nil {
		// secondErr := errors.New("failed in second function")
		// return errors.Join(firstErr, secondErr)
		return fmt.Errorf("secondFunction encountered an error: %w", firstErr)
	}
	return nil
}

func main() {
	err := secondFunction()
	fmt.Println("Original error:", err)

	innerError := errors.Unwrap(err)
	fmt.Println("Unwrapped error:", innerError)

	errs := SomeFunction()

	fmt.Println("Custom error:", errs)

	errs = errors.Unwrap(errs)

	innerError2 := fmt.Errorf("innermost: %w", errs)
	fmt.Println("Unwrapped custom error:", innerError2)
}
