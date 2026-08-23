# Functions & Methods

This module explores Go's function system, including functions, methods, pointers, and closures. Understanding these concepts is essential for writing idiomatic Go code.

## Topics Covered

### 19_functions
**Functions** are first-class citizens in Go with powerful features:
- **Multiple return values**: Return several values simultaneously
- **Named return values**: Pre-declare return variables
- **Variadic functions**: Accept variable number of arguments
- **Function types**: Functions as parameters and return values

**Multiple Returns:**
```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

result, err := divide(10, 2)
if err != nil {
    // Handle error
}
```

**Variadic Functions:**
```go
func sum(numbers ...int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}

sum(1, 2, 3, 4)  // 10
sum([]int{1, 2, 3, 4}...)  // Same result
```

### 20_custom_functions
**User-defined functions** and **function organization**:
- **Function naming**: Clear, descriptive names
- **Function length**: Keep functions focused and concise
- **Parameter grouping**: Related parameters together
- **Documentation**: Clear function documentation

### 21_closures_and_anon_functions
**Closures** capture variables from their surrounding scope:
- **Anonymous functions**: Functions without names
- **Variable capture**: Access to outer scope variables
- **Function literals**: Inline function definitions

**Closure Example:**
```go
func counter() func() int {
    count := 0  // Captured variable
    return func() int {
        count++  // Modifies captured variable
        return count
    }
}

c := counter()
fmt.Println(c())  // 1
fmt.Println(c())  // 2
```

**Common Patterns:**
- **Generators**: Functions that return functions
- **Callbacks**: Functions passed as arguments
- **Resource management**: Cleanup functions

### 22_defer
**Defer** schedules function calls for execution after the surrounding function returns:
- **Resource cleanup**: Files, connections, locks
- **LIFO execution**: Last deferred call executes first
- **Argument evaluation**: Arguments evaluated when defer is called

**Resource Management:**
```go
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()  // Guaranteed to run

    // Process file...
    return nil
}
```

**Multiple Defers:**
```go
func complexOperation() {
    defer cleanup1()
    defer cleanup2()  // Runs first
    defer cleanup3()  // Runs second

    // ... operation ...
    return  // cleanup3, cleanup2, cleanup1 execute
}
```

### 23_pointers
**Pointers** hold memory addresses of values:
- **Address operator**: `&` gets address of value
- **Dereference operator**: `*` gets value at address
- **Zero value**: `nil` for pointer types
- **No pointer arithmetic**: Unlike C/C++

**Basic Usage:**
```go
x := 42
p := &x        // p points to x
fmt.Println(*p) // 42 (dereference)
*p = 50        // Modify x through pointer
fmt.Println(x) // 50
```

**When to Use Pointers:**
- **Modify values**: Change original value in function
- **Large structs**: Avoid copying large data
- **Optional values**: `nil` indicates absence
- **Interface implementation**: Some interfaces require pointers

### 24_pointers_examples
**Practical pointer usage** in real scenarios:
- **Function parameters**: Modify caller's variables
- **Method receivers**: Choose value vs pointer receivers appropriately
- **Data structures**: Linked lists, trees with node references

### 25_pointers_with_maps_and_slices
**Reference semantics** of maps and slices:
- **Slices are references**: Point to underlying arrays
- **Maps are references**: All map operations affect original
- **Copy vs reference**: Understanding when data is shared

**Slice Pointer Behavior:**
```go
s1 := []int{1, 2, 3}
s2 := s1        // s2 references same array
s2[0] = 99      // Modifies both s1 and s2

s3 := make([]int, len(s1))
copy(s3, s1)    // Creates independent copy
```

### 26_types_and_methods
**Method receivers** attach functions to types:
- **Value receivers**: `(t T)` - work with copies
- **Pointer receivers**: `(t *T)` - work with originals
- **Method sets**: Which methods a type supports

**Receiver Types:**
```go
type Counter struct { value int }

// Value receiver - works with copy
func (c Counter) Get() int {
    return c.value
}

// Pointer receiver - modifies original
func (c *Counter) Increment() {
    c.value++
}
```

### 27_value_vs_pointer_method_recievers
**Choosing the right receiver type**:
- **Value receivers**: Immutable operations, small structs
- **Pointer receivers**: Mutable operations, large structs, interface implementation
- **Consistency**: Use pointer receivers if any method needs one

**Interface Implementation:**
```go
type Writer interface {
    Write([]byte) (int, error)
}

// Must use pointer receiver to implement interface
func (f *File) Write(data []byte) (int, error) {
    // Implementation
}
```

### 28_methods_vs_functions
**Methods vs functions** and **object-oriented programming in Go**:
- **Methods**: Functions attached to types
- **Functions**: Standalone functions
- **Composition over inheritance**: Go's OOP approach

**Method Benefits:**
- **Encapsulation**: Data and behavior together
- **Polymorphism**: Different types with same methods
- **Clarity**: `obj.Method()` vs `Method(obj)`