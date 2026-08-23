# Advanced Concurrency Patterns

This module covers **real-world concurrent programming patterns** that solve complex problems. These patterns demonstrate how to build scalable, maintainable concurrent systems using Go's concurrency primitives.

## Topics Covered

### 65_pipelines
**Pipeline pattern** chains goroutines together for data processing:
- **Stages**: Each goroutine performs one transformation
- **Data flow**: Output of one stage becomes input to next
- **Fan-out/fan-in**: Distribute work across multiple goroutines

**Basic Pipeline:**
```go
// Stage 1: Generate numbers
func generator() <-chan int {
    ch := make(chan int)
    go func() {
        for i := 0; i < 10; i++ {
            ch <- i
        }
        close(ch)
    }()
    return ch
}

// Stage 2: Square numbers
func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// Usage
numbers := generator()
squares := square(numbers)
for sq := range squares {
    fmt.Println(sq)
}
```

### 66_pipelines2
**File processing pipelines** apply transformations to large files:
- **Streaming processing**: Process files without loading entirely into memory
- **Error propagation**: Handle errors across pipeline stages
- **Resource cleanup**: Proper file handle management

**File Processing Example:**
```go
func readLines(filename string) <-chan string {
    ch := make(chan string)
    go func() {
        defer close(ch)
        file, err := os.Open(filename)
        if err != nil {
            // Handle error
            return
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            ch <- scanner.Text()
        }
    }()
    return ch
}
```

### 67_66_but_concurrently
**Concurrent pipeline execution** with **worker pools**:
- **Parallel processing**: Multiple workers process items simultaneously
- **Load balancing**: Distribute work across available workers
- **Backpressure**: Handle varying processing speeds

**Worker Pool Pattern:**
```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        results <- process(job)
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    // Start workers
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    // Send jobs
    for j := 1; j <= 9; j++ {
        jobs <- j
    }
    close(jobs)

    // Collect results
    for r := 1; r <= 9; r++ {
        <-results
    }
}
```

### 68_67v2
**Enhanced worker pools** with **dynamic scaling** and **graceful shutdown**:
- **Dynamic workers**: Scale workers based on load
- **Shutdown signaling**: Clean termination of workers
- **Result aggregation**: Collect and combine results

### 69_checking_http
**Concurrent HTTP requests** for improved performance:
- **Batch requests**: Make multiple HTTP calls simultaneously
- **Timeout handling**: Prevent hanging requests
- **Error aggregation**: Collect errors from multiple requests

**Concurrent HTTP Example:**
```go
func fetchURL(url string) (string, error) {
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    // ... process response
}

func main() {
    urls := []string{"url1", "url2", "url3"}
    results := make(chan string, len(urls))

    for _, url := range urls {
        go func(u string) {
            result, err := fetchURL(u)
            if err != nil {
                results <- fmt.Sprintf("Error: %v", err)
            } else {
                results <- result
            }
        }(url)
    }

    // Collect results
    for i := 0; i < len(urls); i++ {
        fmt.Println(<-results)
    }
}
```

### 70_checking_http_fan_in_out
**Fan-in and fan-out patterns** for scalable HTTP processing:
- **Fan-out**: Distribute work to multiple workers
- **Fan-in**: Collect results from multiple workers
- **Load balancing**: Efficient resource utilization

**Fan-Out/Fan-In Pattern:**
```go
// Fan-out: Distribute work
func fanOut(urls []string, workers int) <-chan string {
    ch := make(chan string)
    for i := 0; i < workers; i++ {
        go func(workerID int) {
            for _, url := range urls {
                result := fetchURL(url)
                ch <- result
            }
        }(i)
    }
    return ch
}

// Fan-in: Collect results
func fanIn(inputs ...<-chan string) <-chan string {
    ch := make(chan string)
    for _, input := range inputs {
        go func(in <-chan string) {
            for result := range in {
                ch <- result
            }
        }(input)
    }
    return ch
}
```

### 71_goroutines_and_context
**Context package** for **goroutine lifecycle management**:
- **Cancellation**: Signal goroutines to stop work
- **Timeouts**: Automatic cancellation after time limits
- **Request-scoped values**: Pass request-specific data
- **Tree-structured**: Child contexts inherit from parents

**Context Usage:**
```go
func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            // Cleanup and exit
            return
        default:
            // Do work
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    go worker(ctx, 1)
    go worker(ctx, 2)

    // Workers automatically stop after 5 seconds
    time.Sleep(6 * time.Second)
}
```

**Context Types:**
- `context.Background()`: Root context
- `context.WithCancel()`: Cancellable context
- `context.WithTimeout()`: Time-limited context
- `context.WithValue()`: Context with values