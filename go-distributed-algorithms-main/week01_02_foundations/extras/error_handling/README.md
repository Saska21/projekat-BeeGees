# Error Handling

This module covers Go's unique approach to error handling, which emphasizes **explicit error checking** rather than exceptions. Understanding these patterns is crucial for writing robust Go programs.

## Topics Covered

### 34_error_handling
**Go's error handling philosophy**: Errors are values, not exceptions:
- **Explicit checking**: Must check errors after operations
- **Error interface**: `type error interface { Error() string }`
- **Multiple returns**: Functions return `(result, error)`

**Basic Pattern:**
```go
func doSomething() (int, error) {
    // ... operation that might fail ...
    if failed {
        return 0, errors.New("operation failed")
    }
    return result, nil
}

// Caller must check error
result, err := doSomething()
if err != nil {
    // Handle error
    return err
}
// Use result
```

**Common Error Types:**
- **Sentinel errors**: Predefined error values (`io.EOF`, `sql.ErrNoRows`)
- **Custom errors**: `errors.New()` or custom error types
- **Wrapped errors**: Errors containing other errors

### 35_sentinel_errors_and_wrapping
**Sentinel errors** are predefined error values used for specific conditions:
- **Global variables**: `var ErrNotFound = errors.New("not found")`
- **Comparison**: `if err == io.EOF`
- **Problems**: Can be too generic, lose context

**Error Wrapping** preserves error context through call stacks:
```go
func readFile(filename string) error {
    data, err := os.ReadFile(filename)
    if err != nil {
        return fmt.Errorf("failed to read %s: %w", filename, err)
    }
    // ... process data
}

// Usage
err := readFile("config.json")
if err != nil {
    log.Printf("Configuration error: %v", err)
}
```

**Wrapping Benefits:**
- **Context preservation**: Original error details maintained
- **Debugging**: Full error chain visible
- **Error inspection**: Can check for specific error types

### 36_comparing_errors
**Modern error comparison** using `errors.Is()` and `errors.As()`:
- **`errors.Is()`**: Checks if error matches a target (including wrapped errors)
- **`errors.As()`**: Extracts specific error type from wrapped errors

**Comparison Methods:**
```go
// Old way (doesn't work with wrapped errors)
if err == io.EOF { /* handle */ }

// New way (works with wrapped errors)
if errors.Is(err, io.EOF) { /* handle */ }

// Type assertion
var customErr *MyError
if errors.As(err, &customErr) {
    // Use customErr
}
```

### 37_panic_and_recover
**Panic** and **recover** for exceptional circumstances:
- **Panic**: Immediate termination with stack unwinding
- **Recover**: Catch panics in deferred functions
- **Use sparingly**: Only for truly unexpected errors

**When to Use Panic:**
- **Programming errors**: Array out of bounds, nil pointer dereference
- **Initialization failures**: Cannot proceed if setup fails
- **NOT for expected errors**: File not found, network timeouts

**Recover Pattern:**
```go
func safeFunction() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered from panic: %v", r)
        }
    }()

    // Code that might panic
    riskyOperation()
    return nil
}
```

**Panic Recovery Guidelines:**
- **Recover at boundaries**: HTTP handlers, goroutine entry points
- **Don't overuse**: Most errors should be returned normally
- **Log panics**: Always log recovered panics for debugging
- **Clean up resources**: Ensure proper cleanup even after panic