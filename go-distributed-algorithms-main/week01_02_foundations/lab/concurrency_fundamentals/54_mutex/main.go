package main

import (
	"fmt"
	"sync"
	"time"
)

type Counter struct {
	mu    sync.RWMutex // Read-Write mutex to allow multiple readers but only one writer
	count int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *Counter) Get() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.count
}

func main() {

	counter := &Counter{count: 0}

	wg := sync.WaitGroup{}

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				counter.Inc()
			}
			wg.Done()
		}()
	}

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			fmt.Println("Current Count:", counter.Get())
			time.Sleep(2 * time.Millisecond)
			wg.Done()
		}()
	}

	wg.Wait()
}
