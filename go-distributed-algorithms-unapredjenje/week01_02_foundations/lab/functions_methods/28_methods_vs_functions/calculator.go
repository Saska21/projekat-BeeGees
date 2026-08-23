package main

type Calculator struct {
}

func (c Calculator) Add(a, b int) int {
	return a + b
}

func AddFunction(a, b int) int {
	return a + b
}

func ArithemeticOperation(operation func(int, int) int, a, b int) int {
	return operation(a, b)
}
