package main

import (
	"fmt"
	"sync"
	"time"
)

type CounterMutex struct {
	mu    sync.Mutex
	value int
}

func (c *CounterMutex) Inc() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *CounterMutex) Get() int {
	c.mu.Lock()
	v := c.value
	c.mu.Unlock()
	return v
}

type CounterRWMutex struct {
	mu    sync.RWMutex
	value int
}

func (c *CounterRWMutex) Inc() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *CounterRWMutex) Get() int {
	c.mu.RLock()
	v := c.value
	c.mu.RUnlock()
	return v
}

func main() {
	fmt.Println("=== Concurrency Fundamentals: Exercise Solutions ===")
	task1CounterMutexWaitGroup()
	task2DeadlockAndFix()
	task3ReadHeavyMutexVsRWMutex()
}

func task1CounterMutexWaitGroup() {
	fmt.Println("\n--- Task 1: Counter, Mutex, and WaitGroup ---")

	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("Expected counter = 100")
	fmt.Println("Actual counter   =", counter)
}

func task2DeadlockAndFix() {
	fmt.Println("\n--- Task 2: Create and Fix a Deadlock ---")
	showDeadlock()
	showFixedLockOrdering()
}

func showDeadlock() {
	fmt.Println("Deadlock demo (opposite lock order):")

	var m1, m2 sync.Mutex

	go func() {
		m1.Lock()
		defer m1.Unlock()
		time.Sleep(20 * time.Millisecond)
		m2.Lock()
		defer m2.Unlock()
	}()

	go func() {
		m2.Lock()
		defer m2.Unlock()
		time.Sleep(20 * time.Millisecond)
		m1.Lock()
		defer m1.Unlock()
	}()

	select {
	case <-time.After(120 * time.Millisecond):
		fmt.Println("Likely deadlock observed: goroutines did not finish in time.")
	}
}

func showFixedLockOrdering() {
	fmt.Println("Fixed version (same lock order m1 -> m2):")

	var m1, m2 sync.Mutex
	var wg sync.WaitGroup

	worker := func(id int) {
		defer wg.Done()
		m1.Lock()
		defer m1.Unlock()
		time.Sleep(20 * time.Millisecond)
		m2.Lock()
		defer m2.Unlock()
		fmt.Printf("Worker %d finished without deadlock\n", id)
	}

	wg.Add(2)
	go worker(1)
	go worker(2)
	wg.Wait()
	fmt.Println("Both workers completed (deadlock fixed).")
}

func task3ReadHeavyMutexVsRWMutex() {
	fmt.Println("\n--- Task 3: Read-Heavy Counter with Mutex vs RWMutex ---")

	const readers = 80
	const writers = 4
	const writerOps = 500
	const readerOps = 600

	mCounter := &CounterMutex{}
	rwCounter := &CounterRWMutex{}

	mutexDuration := runReadHeavyWorkloadMutex(mCounter, readers, writers, readerOps, writerOps)
	rwDuration := runReadHeavyWorkloadRWMutex(rwCounter, readers, writers, readerOps, writerOps)

	expected := writers * writerOps
	fmt.Printf("Mutex counter final value:   %d (expected %d), time: %s\n", mCounter.Get(), expected, mutexDuration)
	fmt.Printf("RWMutex counter final value: %d (expected %d), time: %s\n", rwCounter.Get(), expected, rwDuration)
}

func runReadHeavyWorkloadMutex(c *CounterMutex, readers, writers, readerOps, writerOps int) time.Duration {
	start := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < writerOps; j++ {
				c.Inc()
			}
		}()
	}

	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sum := 0
			for j := 0; j < readerOps; j++ {
				sum += c.Get()
			}
			_ = sum
		}()
	}

	wg.Wait()
	return time.Since(start)
}

func runReadHeavyWorkloadRWMutex(c *CounterRWMutex, readers, writers, readerOps, writerOps int) time.Duration {
	start := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < writerOps; j++ {
				c.Inc()
			}
		}()
	}

	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sum := 0
			for j := 0; j < readerOps; j++ {
				sum += c.Get()
			}
			_ = sum
		}()
	}

	wg.Wait()
	return time.Since(start)
}
