# Synchronization Primitives

This module covers advanced synchronization tools that solve specific concurrency problems beyond basic mutexes. These primitives provide more sophisticated coordination between goroutines.

## Topics Covered

### 55_sync.Cond
**Condition Variables** coordinate goroutines waiting for specific conditions:
- **Purpose**: Signal when shared state changes meet certain criteria
- **Components**: Mutex + condition queue + wait/signal operations
- **Use cases**: Producer-consumer patterns, thread pools, event notification

**Key Methods:**
- `Wait()`: Releases mutex and blocks until signaled
- `Signal()`: Wakes one waiting goroutine
- `Broadcast()`: Wakes all waiting goroutines

**Classic Producer-Consumer Example:**
```go
cond := sync.NewCond(&sync.Mutex{})
queue := make([]int, 0, 10)

// Producer
func produce() {
    cond.L.Lock()
    for len(queue) == cap(queue) {
        cond.Wait()  // Wait for space
    }
    queue = append(queue, item)
    cond.Signal()  // Signal consumer
    cond.L.Unlock()
}

// Consumer
func consume() {
    cond.L.Lock()
    for len(queue) == 0 {
        cond.Wait()  // Wait for items
    }
    item := queue[0]
    queue = queue[1:]
    cond.Signal()  // Signal producer
    cond.L.Unlock()
}
```

### 56_sync.Once
**Once** ensures a function executes exactly once, regardless of how many times it's called:
- **Thread-safe singleton initialization**
- **Lazy initialization** - only when first needed
- **No race conditions** - guaranteed single execution

**Common Patterns:**
```go
var (
    instance *Database
    once     sync.Once
)

func GetDatabase() *Database {
    once.Do(func() {
        instance = connectToDatabase()  // Only runs once
    })
    return instance
}
```

**Benefits:**
- **Performance**: Expensive initialization happens once
- **Thread safety**: No race conditions in initialization
- **Lazy loading**: Resources created only when needed

### 57_sync.Pool
**Pool** provides object reuse to reduce garbage collection pressure:
- **Object caching**: Reuse expensive-to-create objects
- **Per-goroutine pools**: Each goroutine has its own pool
- **Automatic cleanup**: GC clears pools when memory pressure increases

**How it works:**
```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 4096)  // Create new buffer if pool empty
    },
}

func processData(data []byte) {
    buf := bufferPool.Get().([]byte)  // Get from pool
    defer bufferPool.Put(buf)         // Return to pool

    // Use buf...
}
```

**When to use:**
- **Buffers**: `bytes.Buffer`, network buffers
- **Workers**: Reusable worker structs
- **Encoders/Decoders**: `json.Encoder`, `xml.Decoder`
- **Any expensive object** that can be reset and reused

**Performance Impact:**
- **Reduced allocations**: Reuse existing objects
- **Lower GC pressure**: Fewer objects for garbage collector
- **Better cache locality**: Objects stay in memory longer