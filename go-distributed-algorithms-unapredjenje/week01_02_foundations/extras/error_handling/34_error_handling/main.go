package main

import (
	"fmt"
)

type CustomError struct {
	Code    int
	Message string
}

func (c CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", c.Code, c.Message)
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		//return 0, errors.New("cannot divide by zero")
		return 0, CustomError{Code: 400, Message: fmt.Sprintf("cannot divide %d by zero", a)}
	}
	return a / b, nil
}

func main() {
	result, err := Divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Result:", result)
}
