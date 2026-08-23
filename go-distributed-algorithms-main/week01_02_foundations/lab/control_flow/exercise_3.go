package main

import "fmt"

func main() {
	fmt.Println("=== Exercise 3 Solutions (Control Flow) ===")
	task1IfInitializer()
	task2ShadowingVsReassignment()
	task3LoopControl()
	task4LabeledBreak()
	task5Goto()
}

func isEven(n int) bool {
	return n%2 == 0
}

func task1IfInitializer() {
	fmt.Println("\n--- Task 1: If Initializer ---")

	n := 11
	if even := isEven(n); even {
		fmt.Printf("%d is even\n", n)
	} else {
		fmt.Printf("%d is odd\n", n)
	}

	// This fails because 'even' exists only in the if/else scope.
	// Uncomment to see compile error: undefined: even
	// fmt.Println(even)
	fmt.Println("Reason: variable declared in if initializer is scoped only to that if block.")
}

func task2ShadowingVsReassignment() {
	fmt.Println("\n--- Task 2: Shadowing vs Reassignment ---")

	a := 10
	fmt.Println("Before block, a =", a)

	{
		a := 20
		fmt.Println("Inside block after shadowing (a := 20), a =", a)

		a = 30
		fmt.Println("Inside block after reassignment (a = 30), a =", a)
	}

	fmt.Println("After block, outer a =", a)
}

func task3LoopControl() {
	fmt.Println("\n--- Task 3: Loop Control (continue + break) ---")

	for i := 0; i <= 10; i++ {
		if i%2 == 0 {
			continue
		}
		if i > 7 {
			break
		}
		fmt.Println(i)
	}
}

func task4LabeledBreak() {
	fmt.Println("\n--- Task 4: Labeled Break ---")
	fmt.Println("Normal break (only exits inner loop):")

	for i := 0; i <= 3; i++ {
		for j := 0; j <= 3; j++ {
			if i == 2 && j == 1 {
				fmt.Println("normal break at", i, j)
				break
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
		fmt.Println()
	}

	fmt.Println("Labeled break (exits both loops):")
outer:
	for i := 0; i <= 3; i++ {
		for j := 0; j <= 3; j++ {
			if i == 2 && j == 1 {
				fmt.Println("labeled break at", i, j)
				break outer
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
		fmt.Println()
	}
}

func task5Goto() {
	fmt.Println("\n--- Task 5: goto ---")

	for i := 0; i <= 5; i++ {
		if i == 3 {
			goto end
		}
		fmt.Println("i =", i)
	}

end:
	fmt.Println("Skipped to end")
}
