# GAction

A Go library that provides C#-style events and delegates with full generic type support. GAction enables type-safe event handling and callback management patterns commonly found in C# and other .NET languages.

## Features

- **Type-safe delegates** with 0-4 parameters
- **Event system** with subscribe/unsubscribe functionality
- **Thread-safe** implementation using RW mutexes
- **Generic support** for any type combinations
- **Memory-efficient** subscription management
- **Interface-based design** for better abstraction

## Installation

```bash
go get github.com/DilemaFixer/GAction
```

## Quick Start

### Basic Delegate Usage

```go
package main

import (
    "fmt"
    "github.com/DilemaFixer/GAction"
)

func main() {
    // Create a delegate that takes no parameters
    delegate := GAction.NewDelegate0()
    delegate.Set(func() {
        fmt.Println("Hello from delegate!")
    })
    
    delegate.Invoke() // Output: Hello from delegate!
}
```

### Basic Event Usage

```go
package main

import (
    "fmt"
    "github.com/DilemaFixer/GAction"
)

func main() {
    // Create an event with one string parameter
    onMessage, emitMessage := GAction.NewEvent1[string]()
    
    // Subscribe multiple handlers
    onMessage.Subscribe(func(msg string) {
        fmt.Println("Handler 1:", msg)
    })
    
    onMessage.Subscribe(func(msg string) {
        fmt.Println("Handler 2:", msg)
    })
    
    // Emit the event
    emitMessage.Emit("Hello World!")
    // Output:
    // Handler 1: Hello World!
    // Handler 2: Hello World!
}
```

## Core Components

### Delegates

Delegates are type-safe function pointers that can be set and invoked. They support 0-4 parameters:

- `Delegate0` - No parameters
- `Delegate1[T]` - One parameter of type T
- `Delegate2[A, B]` - Two parameters of types A and B
- `Delegate3[A, B, C]` - Three parameters
- `Delegate4[A, B, C, D]` - Four parameters

#### Delegate Example

```go
// Delegate with one integer parameter
delegate := GAction.NewDelegate1[int]()
delegate.Set(func(value int) {
    fmt.Printf("Received: %d\n", value)
})

delegate.Invoke(42) // Output: Received: 42
```

### Events

Events provide a publish-subscribe pattern where multiple subscribers can listen to a single event. Like delegates, they support 0-4 parameters.

#### Event Types

- `Event0` / `Emitter0` - No parameters
- `Event1[T]` / `Emitter1[T]` - One parameter
- `Event2[A, B]` / `Emitter2[A, B]` - Two parameters
- `Event3[A, B, C]` / `Emitter3[A, B, C]` - Three parameters
- `Event4[A, B, C, D]` / `Emitter4[A, B, C, D]` - Four parameters

#### Multi-parameter Event Example

```go
// Event with two parameters: string tag and int code
onLog, emitLog := GAction.NewEvent2[string, int]()

onLog.Subscribe(func(tag string, code int) {
    fmt.Printf("[%s] Status Code: %d\n", tag, code)
})

emitLog.Emit("INFO", 200)  // Output: [INFO] Status Code: 200
emitLog.Emit("ERROR", 404) // Output: [ERROR] Status Code: 404
```

### Subscriptions and Unsubscribing

All event subscriptions return a `Subscription` interface that allows you to unsubscribe:

```go
onUpdate, emitUpdate := GAction.NewEvent1[string]()

// Subscribe and keep the subscription reference
subscription := onUpdate.Subscribe(func(msg string) {
    fmt.Println("Update:", msg)
})

emitUpdate.Emit("First update")  // Will print
subscription.Unsubscribe()       // Remove this subscriber
emitUpdate.Emit("Second update") // Won't print from unsubscribed handler
```

## Advanced Usage

### Interface-Based Design

Use interfaces to expose only the functionality you want:

```go
type Worker struct {
    OnComplete GAction.Invoker1[int] // Read-only interface
}

func NewWorker() *Worker {
    delegate := GAction.NewDelegate1[int]()
    delegate.Set(func(result int) {
        fmt.Printf("Work completed with result: %d\n", result)
    })
    
    return &Worker{
        OnComplete: delegate, // Expose only the Invoker interface
    }
}

func main() {
    worker := NewWorker()
    worker.OnComplete.Invoke(100) // Can invoke, but not change the handler
}
```

### Event Chaining

Events can trigger other events, creating chains of behavior:

```go
onStart, emitStart := GAction.NewEvent0()
onProcess, emitProcess := GAction.NewEvent0()
onComplete, emitComplete := GAction.NewEvent0()

// Chain events: start -> process -> complete
onStart.Subscribe(func() {
    fmt.Println("Starting...")
    emitProcess.Emit()
})

onProcess.Subscribe(func() {
    fmt.Println("Processing...")
    emitComplete.Emit()
})

onComplete.Subscribe(func() {
    fmt.Println("Completed!")
})

emitStart.Emit()
// Output:
// Starting...
// Processing...
// Completed!
```

### Service Pattern with Events

```go
type Service struct {
    OnReady  GAction.Event0
    OnError  GAction.Event1[error]
}

func NewService() (*Service, GAction.Emitter0, GAction.Emitter1[error]) {
    readyEvent, readyEmitter := GAction.NewEvent0()
    errorEvent, errorEmitter := GAction.NewEvent1[error]()
    
    return &Service{
        OnReady: readyEvent,
        OnError: errorEvent,
    }, readyEmitter, errorEmitter
}

func main() {
    service, ready, errorEmitter := NewService()
    
    service.OnReady.Subscribe(func() {
        fmt.Println("Service is ready!")
    })
    
    service.OnError.Subscribe(func(err error) {
        fmt.Printf("Service error: %v\n", err)
    })
    
    ready.Emit() // Trigger ready event
}
```

## Utility Functions

### Combining Actions

You can combine multiple action functions into a single action:

```go
action1 := func(msg string) { fmt.Println("Action 1:", msg) }
action2 := func(msg string) { fmt.Println("Action 2:", msg) }
action3 := func(msg string) { fmt.Println("Action 3:", msg) }

combined := GAction.Combine1(action1, action2, action3)
combined("Hello") 
// Output:
// Action 1: Hello
// Action 2: Hello  
// Action 3: Hello
```

## Thread Safety

All components in GAction are thread-safe and can be safely used across multiple goroutines:

```go
onData, emitData := GAction.NewEvent1[int]()

onData.Subscribe(func(value int) {
    fmt.Printf("Goroutine received: %d\n", value)
})

// Safe to emit from multiple goroutines
go func() { emitData.Emit(1) }()
go func() { emitData.Emit(2) }()
go func() { emitData.Emit(3) }()
```

## Best Practices

1. **Use interfaces** to expose only needed functionality (e.g., `Invoker1[T]` instead of `*Delegate1[T]`)
2. **Always check subscriptions** - unsubscribe when no longer needed to prevent memory leaks
3. **Handle nil actions** - the library safely handles nil function pointers
4. **Separate concerns** - use separate events for different types of notifications
5. **Consider event ordering** - subscribers are called in the order they were added

## API Reference

### Types

- `Action0` through `Action4[...]` - Function type aliases
- `Delegate0` through `Delegate4[...]` - Mutable function holders
- `Event0` through `Event4[...]` - Event subscription interfaces
- `Emitter0` through `Emitter4[...]` - Event emission interfaces
- `Invoker0` through `Invoker4[...]` - Read-only invocation interfaces
- `Subscription` - Unsubscribe interface

### Key Methods

- `NewDelegate0()` through `NewDelegate4[...]()` - Create delegates
- `NewEvent0()` through `NewEvent4[...]()` - Create events (returns Event and Emitter)
- `Set(func)` - Set delegate function
- `Invoke(...)` - Call delegate or invoke subscribers
- `Subscribe(func) Subscription` - Subscribe to event
- `Emit(...)` - Emit event to all subscribers
- `Unsubscribe()` - Remove subscription
- `HasSubscribers() bool` - Check if event has active subscribers
