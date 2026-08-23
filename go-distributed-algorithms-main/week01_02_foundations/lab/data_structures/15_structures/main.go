package main

import "fmt"

func main() {
	// Usually outside of main
	type student struct {
		firstName string
		lastName  string
		age       int
		subjects  []string
	}

	var student1 student
	fmt.Printf("%+v", student1)

	// Initializing struct
	// Explicit field names
	student1 = student{
		firstName: "John",
		lastName:  "Doe",
		age:       21,
		subjects:  []string{"Math", "Science"},
	}
	fmt.Printf("\n%+v\n", student1)

	// Without field names (order matters)
	student2 := student{"Jane", "Smith", 22, []string{"History", "Art"}}
	fmt.Printf("\n%+v\n", student2)

	// Reading and writing fields
	fmt.Println("First Name of student1:", student1.firstName)

	student1.subjects = append(student1.subjects, "English")
	fmt.Printf("Updated subjects of student1: %+v\n", student1.subjects)

	// Anonymous struct
	guardian := struct {
		firstName string
		lastName  string
	}{
		firstName: "Alice",
		lastName:  "Johnson",
	}
	fmt.Printf("Guardian: %+v\n", guardian)
}
