# Concurrency Fundamentals

This module introduces Go's powerful concurrency model, which allows programs to execute multiple tasks simultaneously. Unlike traditional threading, Go uses **goroutines** - lightweight threads managed by the Go runtime.

## Topics Covered

### 49_concurrency
**Goroutines** are the foundation of Go's concurrency model:
- **Lightweight threads**: Goroutines are managed by the Go runtime, not the OS
- **Low memory footprint**: Initial stack of ~2KB vs ~2MB for OS threads
- **Fast creation**: Thousands of goroutines can be created without performance issues
- **Automatic scheduling**: Go runtime handles goroutine scheduling across CPU cores

**Key Concepts:**
- `go` keyword launches functions concurrently
- Main function runs as a goroutine
- Program exits when main goroutine finishes (regardless of other goroutines)

### 50_deadlock
**Deadlock** occurs when goroutines wait indefinitely for resources held by each other:
- **Four conditions for deadlock**: Mutual exclusion, hold and wait, no preemption, circular wait
- **Common causes**: Nested lock acquisition, forgotten unlocks, circular dependencies
- **Prevention**: Consistent lock ordering, timeout mechanisms, deadlock detection

**Example deadlock scenario:**
```go
// Goroutine 1: Lock A then B
mutexA.Lock()
mutexB.Lock() // Waits forever if Goroutine 2 has B

// Goroutine 2: Lock B then A
mutexB.Lock()
mutexA.Lock() // Waits forever if Goroutine 1 has A
```

### 51_livelock
**Livelock** is when goroutines are active but make no progress:
- **Difference from deadlock**: Goroutines are running, not blocked
- **Symptoms**: High CPU usage, no forward progress, endless retries
- **Common cause**: Overly aggressive retry logic without backoff

**Example**: Two people trying to pass each other in a hallway, both stepping the same way repeatedly.

### 52_starvation
**Starvation** happens when a goroutine cannot access a resource because others monopolize it:
- **Priority inversion**: Low-priority tasks starve high-priority ones
- **Fairness issues**: Some goroutines never get CPU time or resource access
- **Mutex unfairness**: Later-arriving goroutines may wait indefinitely

**Real-world example**: A busy web server where background cleanup tasks never run because request handlers consume all resources.

### 53_sync.WaitGroup
**WaitGroup** synchronizes goroutine completion:
- **Counter-based synchronization**: Tracks number of active goroutines
- **Methods**: `Add()`, `Done()`, `Wait()`
- **Use case**: Ensure main function waits for all goroutines to complete

**Pattern:**
```go
var wg sync.WaitGroup
wg.Add(3)  // Expect 3 goroutines

go func() { defer wg.Done(); /* work */ }()
go func() { defer wg.Done(); /* work */ }()
go func() { defer wg.Done(); /* work */ }()

wg.Wait()  // Blocks until counter reaches 0
```

### 54_mutex
**Mutex** (Mutual Exclusion) prevents concurrent access to shared resources:
- **Critical sections**: Code that must run atomically
- **Lock types**: `Lock()` (exclusive) and `RLock()` (read-shared)
- **Deadlock prevention**: Always unlock in defer statements

**Read-Write Mutex (`sync.RWMutex`)**:
- Multiple readers can hold the lock simultaneously
- Only one writer can hold the lock (exclusive)
- Writers block all readers and other writers