# Channels

This module explores Go's **channel-based concurrency model** - the primary mechanism for goroutine communication. Channels provide **type-safe, synchronized communication** between goroutines, eliminating the need for explicit locks in many cases.

## Topics Covered

### 58_channels
**Channels** are typed conduits for sending and receiving values between goroutines:
- **Typed communication**: `chan int`, `chan string`, etc.
- **Synchronous by default**: Send/receive operations block until counterpart is ready
- **Thread-safe**: No race conditions when using channels correctly

**Basic Operations:**
```go
ch := make(chan int)  // Create channel

go func() {
    ch <- 42  // Send value (blocks until received)
}()

value := <-ch  // Receive value (blocks until sent)
```

**Channel Types:**
- **Bidirectional**: `chan T` - can send and receive
- **Send-only**: `chan<- T` - can only send
- **Receive-only**: `<-chan T` - can only receive

### 59_channels2
**Advanced channel patterns** and **buffered channels**:
- **Buffered channels**: `make(chan int, 10)` - holds N values before blocking
- **Non-blocking operations**: Send/receive without blocking
- **Channel closing**: `close(ch)` signals no more values will be sent

**Buffered vs Unbuffered:**
```go
// Unbuffered (synchronous)
ch1 := make(chan int)        // Blocks immediately
ch1 <- 1                     // Blocks until received

// Buffered (asynchronous)
ch2 := make(chan int, 2)     // Can hold 2 values
ch2 <- 1                     // Doesn't block
ch2 <- 2                     // Doesn't block
ch2 <- 3                     // Blocks (buffer full)
```

**Closing Channels:**
```go
ch := make(chan int, 10)
// ... send values ...
close(ch)  // No more sends allowed

// Receiving from closed channel:
value, ok := <-ch
if !ok {
    // Channel is closed, no more values
}
```

### 60_select
**Select** multiplexes multiple channel operations, allowing goroutines to wait on multiple channels simultaneously:
- **Non-blocking choice**: Proceed with whichever operation is ready first
- **Timeout handling**: Implement timeouts with `time.After()`
- **Default case**: Non-blocking operations

**Basic Select:**
```go
select {
case msg := <-ch1:
    fmt.Println("Received from ch1:", msg)
case msg := <-ch2:
    fmt.Println("Received from ch2:", msg)
case ch3 <- "hello":
    fmt.Println("Sent to ch3")
}
```

**Timeout Pattern:**
```go
select {
case result := <-ch:
    return result
case <-time.After(5 * time.Second):
    return errors.New("timeout")
}
```

**Non-blocking Select:**
```go
select {
case x := <-ch:
    fmt.Println("Received:", x)
default:
    fmt.Println("No value ready")
}
```

### 61_for_select
**Range over channels** and **select in loops** for continuous processing:
- **Range loop**: `for value := range ch` - receives until channel closes
- **Event loop**: Continuous select for handling multiple inputs
- **Graceful shutdown**: Proper channel closing patterns

**Range Pattern:**
```go
ch := make(chan int, 10)
// ... send values ...
close(ch)  // Signal end

for value := range ch {
    fmt.Println("Received:", value)
}
// Loop exits when channel closes
```

**Event Loop Pattern:**
```go
for {
    select {
    case req := <-requests:
        processRequest(req)
    case <-shutdown:
        cleanup()
        return
    }
}
```

**Shutdown Pattern:**
```go
done := make(chan struct{})
go func() {
    defer close(done)
    // ... work ...
}()

// Wait for completion
<-done
```