// week03_04_basic_abstractions - Teaching Simulator for Distributed Algorithms
//
// This program demonstrates the basic abstractions from Chapter 2 of
// Cachin, Guerraoui, Rodrigues: "Reliable and Secure Distributed Programming"
//
// It runs three self-contained examples:
//  1. Ping-Pong:  Basic message passing with perfect links
//  2. Counter:    Message loss with fair-loss links
//  3. Tamper:     Message authentication and forgery detection
//
// Usage:
//
//	go run .
//
// KEY CONSTRAINT: No goroutines, no channels, no networking.
// Everything is single-threaded and deterministic given a fixed random seed.
package main

import (
	"fmt"

	"github.com/danilokacanski/da/week03_04_basic_abstractions/examples"
)

func main() {
	fmt.Println("================================================================")
	fmt.Println("  DA Week 2-3: Basic Abstractions")
	fmt.Println("  Based on: Cachin, Guerraoui, Rodrigues (Chapter 2)")
	fmt.Println("================================================================")
	fmt.Println()

	// examples.RunPingPong()
	// fmt.Println()

	// examples.RunCounter()
	// fmt.Println()

	examples.RunTamper()
	fmt.Println()

	// fmt.Println("================================================================")
	// fmt.Println("  All examples complete. Review the traces above.")
	// fmt.Println("================================================================")
}
