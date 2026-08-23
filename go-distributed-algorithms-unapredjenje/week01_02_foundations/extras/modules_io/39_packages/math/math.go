package math

// PI (exported) because it starts with an uppercase letter
const PI = 3.14159

// Add (exported) because it starts with an uppercase letter
func Add(a, b int) int {
	return a + b
}

// subtract (unexported) because it starts with a lowercase letter
func subtract(a, b int) int {
	return a - b
}
