package main

import "fmt"

type student struct {
	name string
	age  int
}

func main() {
	fmt.Println("=== Exercise 2 Solutions ===")
	task1ArraysVsSlices()
	task2SliceCapacityAndAppend()
	task3MapsAndStructs()
	bonusSliceRippleEffect()
}

func task1ArraysVsSlices() {
	fmt.Println("\n--- Task 1: Arrays vs Slices ---")

	arr := [4]int{1, 2, 3, 4}
	s := arr[1:3]

	fmt.Println("Initial array:", arr)
	fmt.Println("Initial slice:", s)

	s[0] = 99
	fmt.Println("After modifying slice:")
	fmt.Println("Array:", arr)
	fmt.Println("Slice:", s)

	sCopy := make([]int, len(s))
	copy(sCopy, s)
	sCopy[0] = 777

	fmt.Println("After modifying copied slice:")
	fmt.Println("Array (unchanged by copy edit):", arr)
	fmt.Println("Original slice:", s)
	fmt.Println("Copied slice:", sCopy)
}

func task2SliceCapacityAndAppend() {
	fmt.Println("\n--- Task 2: Slice Capacity and append ---")

	s := make([]int, 2, 3)
	s[0], s[1] = 10, 20

	beforeAppend := s

	fmt.Printf("Before append: s=%v len=%d cap=%d\n", s, len(s), cap(s))
	fmt.Printf("Before append (alias): beforeAppend=%v len=%d cap=%d\n", beforeAppend, len(beforeAppend), cap(beforeAppend))

	s = append(s, 30)
	fmt.Printf("After 1st append: s=%v len=%d cap=%d\n", s, len(s), cap(s))

	s = append(s, 40)
	fmt.Printf("After 2nd append: s=%v len=%d cap=%d\n", s, len(s), cap(s))

	beforeAppend[0] = 111
	s[1] = 222

	fmt.Println("After modifications:")
	fmt.Println("beforeAppend:", beforeAppend)
	fmt.Println("s:", s)
	fmt.Println("Observation: after reallocation, edits in s no longer ripple to beforeAppend.")
}

func task3MapsAndStructs() {
	fmt.Println("\n--- Task 3: Maps and Structs ---")

	class := map[string]student{}

	class["s1"] = student{name: "Alice", age: 20}
	class["s2"] = student{name: "Bob", age: 22}
	fmt.Println("After insert:", class)

	class["s2"] = student{name: "Bob", age: 23}
	fmt.Println("After update s2:", class)

	delete(class, "s1")
	fmt.Println("After delete s1:", class)

	v, ok := class["does_not_exist"]
	fmt.Println("Read missing key -> value:", v, "exists:", ok)
}

func bonusSliceRippleEffect() {
	fmt.Println("\n--- Bonus: Slice Ripple Effect ---")

	s := []int{1, 2, 3, 4, 5}
	a := s[1:4]
	b := s[2:5]

	a[1] = 999
	a[0] = 777

	fmt.Println("Original slice s:", s)
	fmt.Println("Slice a:", a)
	fmt.Println("Slice b:", b)
	fmt.Println("Explanation: a and b overlap the same backing array, so writes in a appear in s and b.")
}
