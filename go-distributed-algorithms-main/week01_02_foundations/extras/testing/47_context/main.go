package main

import (
	"context"
	"fmt"
	"time"
)

func SimulateLongRunningTask(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task cancelled")
			return
		default:
			fmt.Println("Working...")
			time.Sleep(1 * time.Second)
		}
	}
}

func ProcessOrder(ctx context.Context, orderID int) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		fmt.Println("No user ID found in context")
		return
	}
	fmt.Printf("Processing order %d for user %d\n", orderID, userID)
}

func GetUserAndProcessOrder(orderID int) {
	ctx := context.WithValue(context.Background(), "user_id", 12345)
	ProcessOrder(ctx, orderID)
}

func main() {
	// // Example 1:

	// // Create a context with a timeout of 5 seconds
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel() // Ensure cancellation when main function exits

	// // Launch the long-running task in a separate goroutine
	// go SimulateLongRunningTask(ctx)

	// // Simulate main thread doing other work for a while
	// time.Sleep(2 * time.Second)

	// // Cancel the context after some time
	// fmt.Println("Cancelling context...")
	// cancel()

	// // Wait for a moment to allow the goroutine to print the cancellation message
	// time.Sleep(1 * time.Second)

	// Example 2:
	GetUserAndProcessOrder(101)
}
