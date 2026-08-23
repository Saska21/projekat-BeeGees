package main

import "fmt"

type Account struct {
	Balance int
}

func (a Account) depositValue(amount int) {
	a.Balance += amount
}

func (a *Account) depositPointer(amount int) {
	a.Balance += amount
}

func makeBankAccount(initial int) func(int) int {
	balance := initial
	return func(amount int) int {
		balance += amount
		return balance
	}
}

func testDeferClosure() {
	x := 10
	defer fmt.Println("defer direct:", x)
	defer func() {
		fmt.Println("defer closure:", x)
	}()
	x = 20
	fmt.Println("inside test, x changed to:", x)
}

func main() {
	fmt.Println("=== Functions/Methods Exercise Solutions ===")

	fmt.Println("\n--- Task 1: Closures as State ---")
	acc1 := makeBankAccount(100)
	acc2 := makeBankAccount(500)
	fmt.Println("acc1 +50  ->", acc1(50))
	fmt.Println("acc2 +100 ->", acc2(100))
	fmt.Println("acc1 -30  ->", acc1(-30))
	fmt.Println("acc2 +20  ->", acc2(20))
	fmt.Println("acc1 +10  ->", acc1(10))

	fmt.Println("\n--- Task 2: Copy vs Original ---")
	account := Account{Balance: 100}
	fmt.Println("initial balance:", account.Balance)
	account.depositValue(50)
	fmt.Println("after depositValue(50):", account.Balance, "(unchanged, value receiver)")
	account.depositPointer(50)
	fmt.Println("after depositPointer(50):", account.Balance, "(updated, pointer receiver)")

	fmt.Println("\n--- Task 3: Defer and Closures ---")
	testDeferClosure()
}
