package main

import "fmt"

type CustomError struct {
	Message string
	Wrapped error
}

func (e CustomError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Wrapped)
}

func (e CustomError) Unwrap() error {
	return e.Wrapped
}

func SomeFunction() error {
	return CustomError{
		Message: "an error occurred in SomeFunction",
		Wrapped: fmt.Errorf("original error"),
	}
}
