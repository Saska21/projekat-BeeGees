package main

import (
	"errors"
	"fmt"
)

type student struct {
	firstName string
	lastName  string
}

func main() {
	var a int = 10

	increment(&a)
	fmt.Println(a)

	s := student{firstName: "John", lastName: "Doe"}
	previousLastName, err := updateLastName(&s, "Random")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Previous last name:", previousLastName)
	}
	fmt.Printf("Updated student: %+v\n", s)
}

func increment(x *int) {
	(*x)++
}

// updateLastName takes a pointer to a student struct and updates the last name field.
func updateLastName(s *student, newLastName string) (*string, error) {
	if newLastName == "" {
		return nil, errors.New("invalid last name")
	}
	previous := s.lastName
	s.lastName = newLastName
	return &previous, nil
}
