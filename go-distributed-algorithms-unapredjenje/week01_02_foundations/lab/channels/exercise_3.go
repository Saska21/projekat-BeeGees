package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Channels Exercise Solutions ===")
	task1ProducerConsumerSum(10)
	task2BroadcastWithClose()
	task3ForSelectWithTimeout()
}

func task1ProducerConsumerSum(n int) {
	fmt.Println("\n--- Task 1: Producer, Consumer, Sum ---")

	ch := make(chan int)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		sum := 0
		for v := range ch {
			sum += v
		}
		fmt.Printf("Final sum 1..%d = %d\n", n, sum)
	}()

	go func() {
		for i := 1; i <= n; i++ {
			ch <- i
		}
		close(ch)
	}()

	wg.Wait()
}

func task2BroadcastWithClose() {
	fmt.Println("\n--- Task 2: Broadcast with close(signal) ---")

	signal := make(chan struct{})
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-signal
			fmt.Printf("worker %d started\n", id)
		}(i)
	}

	time.Sleep(100 * time.Millisecond)
	close(signal)
	wg.Wait()
}

func task3ForSelectWithTimeout() {
	fmt.Println("\n--- Task 3: for-select with Timeout ---")

	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for i := 1; i <= 4; i++ {
			time.Sleep(200 * time.Millisecond)
			ch1 <- i
		}
	}()

	go func() {
		for i := 100; i <= 103; i++ {
			time.Sleep(300 * time.Millisecond)
			ch2 <- i
		}
	}()

	for {
		select {
		case v := <-ch1:
			fmt.Printf("received %d from ch1\n", v)
		case v := <-ch2:
			fmt.Printf("received %d from ch2\n", v)
		case <-time.After(900 * time.Millisecond):
			fmt.Println("graceful stop")
			return
		}
	}
}
