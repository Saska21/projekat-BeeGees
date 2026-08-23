// week05_06_consensus — Flooding Consensus Teaching Simulator
//
// Implements Algorithm 5.1 (Flooding Consensus) from:
//
//	Cachin, Guerraoui, Rodrigues: "Reliable and Secure Distributed Programming"
//	Chapter 5.1.2
//
// Reuses the concurrent runtime and link stack from week03_04_parallel.
//
// Usage:
//
//	go run .
package main

import (
	"fmt"

	"github.com/danilokacanski/da/week05_06_consensus/examples"
)

func main() {
	fmt.Println("================================================================")
	fmt.Println("  DA Week 5-6: Consensus — Flooding Consensus (Algorithm 5.1)")
	fmt.Println("  Based on: Cachin, Guerraoui, Rodrigues (Chapter 5)")
	fmt.Println("  Runtime: goroutines + channels (week03_04_parallel)")
	fmt.Println("================================================================")
	fmt.Println()

	examples.RunNoFailure()
	fmt.Println()

	// examples.RunOneCrash()
	// fmt.Println()

	// examples.RunDelayedDecision()
	// fmt.Println()

	fmt.Println()
	fmt.Println("================================================================")
	fmt.Println("  All examples complete. Review traces and checks above.")
	fmt.Println("================================================================")
}
