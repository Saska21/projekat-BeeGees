# Advanced Types

This module explores Go's type system beyond basic types, focusing on **interfaces**, **composition**, and **type relationships**. These features enable flexible, maintainable code through polymorphism and code reuse.

## Topics Covered

### 29_deriving_types_and_iota
**Type derivation** creates new types based on existing ones:
- **Type aliases**: `type MyInt = int` (same type)
- **Type definitions**: `type MyInt int` (new type)
- **Method sets**: Different types have different method sets

**iota** for enumerated constants:
```go
type Status int
const (
    Pending Status = iota  // 0
    Running                // 1
    Completed              // 2
    Failed                 // 3
)
```

**Type Safety Benefits:**
- **Intent clarity**: `type UserID int` vs `int`
- **Method attachment**: Only `UserID` gets user-specific methods
- **Type conversion**: Explicit conversion required

### 30_composition_and_promotion
**Composition** builds complex types from simpler ones (Go's alternative to inheritance):
- **Embedding**: Include one struct inside another
- **Field promotion**: Embedded fields accessible at top level
- **Method promotion**: Embedded type methods available on outer type

**Embedding Example:**
```go
type Engine struct {
    Power int
    Type  string
}

func (e Engine) Start() { fmt.Println("Engine started") }

type Car struct {
    Make  string
    Model string
    Engine  // Embedded (no field name)
}

car := Car{Make: "Toyota", Engine: Engine{Power: 200}}
car.Start()      // Promoted method
car.Power        // Promoted field
```

**Composition vs Inheritance:**
- **Flexible**: Can embed multiple types
- **Explicit**: No automatic method overriding
- **Clear**: No confusion about which method runs

### 31_iinterfaces
**Interfaces** define behavior contracts without implementation:
- **Implicit implementation**: Types implement interfaces by having required methods
- **Polymorphism**: Same interface, different implementations
- **Decoupling**: Depend on behavior, not concrete types

**Interface Definition:**
```go
type Writer interface {
    Write([]byte) (int, error)
}

type Reader interface {
    Read([]byte) (int, error)
}

type ReadWriter interface {
    Reader
    Writer
}
```

**Implementation:**
```go
type File struct { /* fields */ }

func (f *File) Write(data []byte) (int, error) {
    // Implementation
}

func (f *File) Read(data []byte) (int, error) {
    // Implementation
}

// File automatically implements ReadWriter
```

### 32_empty_interfaces_type_assertions_switch
**Empty interface** (`interface{}`) accepts any type:
- **Type erasure**: Can hold any value
- **Type assertions**: Check and extract actual type
- **Type switches**: Handle different types in switch statement

**Type Assertions:**
```go
var i interface{} = "hello"

s := i.(string)        // Panic if not string
s, ok := i.(string)    // Safe assertion

if s, ok := i.(string); ok {
    fmt.Println("String:", s)
}
```

**Type Switch:**
```go
switch v := i.(type) {
case int:
    fmt.Printf("Integer: %d\n", v)
case string:
    fmt.Printf("String: %s\n", v)
default:
    fmt.Printf("Unknown type: %T\n", v)
}
```

### 33_nil_interface
**Nil interfaces** and **interface values**:
- **Interface value**: `(type, value)` pair
- **Nil interface**: `nil` type and `nil` value
- **Non-nil interface**: Concrete type with possibly nil value

**Interface Representation:**
```go
var i interface{}  // (nil, nil)

i = (*int)(nil)    // (int*, nil) - non-nil interface with nil value
i = 42            // (int, 42)   - non-nil interface with value

if i == nil {      // Compares both type and value
    // Only true for (nil, nil)
}
```

**Common Pitfalls:**
```go
type error interface {
    Error() string
}

var err error = (*MyError)(nil)  // Non-nil interface
if err != nil {                  // True!
    // err is not nil, even though the value is nil
}
```