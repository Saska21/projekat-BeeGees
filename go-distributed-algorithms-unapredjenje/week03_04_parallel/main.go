// week03_04_parallel - Teaching Simulator for Distributed Algorithms
//
// PARALLEL VERSION: Uses goroutines and channels instead of a
// single-threaded step loop. Each process runs in its own goroutine.
//
// This program demonstrates the basic abstractions from Chapter 2 of
// Cachin, Guerraoui, Rodrigues: "Reliable and Secure Distributed Programming"
//
// It runs three self-contained examples:
//  1. Ping-Pong:  Basic message passing with perfect links (goroutines)
//  2. Counter:    Message loss with fair-loss links (WrapHandler adapter)
//  3. Tamper:     Message authentication and forgery detection (goroutines)
//
// Usage:
//
//	go run .
package main

import (
	"fmt"

	"github.com/danilokacanski/da/week03_04_parallel/examples"
)

func main() {
	fmt.Println("================================================================")
	fmt.Println("  DA Week 2-3: Basic Abstractions (PARALLEL VERSION)")
	fmt.Println("  Based on: Cachin, Guerraoui, Rodrigues (Chapter 2)")
	fmt.Println("  Using: Goroutines, Channels, sync.Mutex")
	fmt.Println("================================================================")
	fmt.Println()

	examples.RunPingPong()
	fmt.Println()

	// examples.RunCounter()
	// fmt.Println()

	// examples.RunTamper()
	// fmt.Println()

	// fmt.Println("================================================================")
	// fmt.Println("  All examples complete. Review the traces above.")
	// fmt.Println("================================================================")
}
