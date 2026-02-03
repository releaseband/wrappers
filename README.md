# Key Wrapper Library

A Go library for automatic key sharding with dynamic shard count adjustment.

## Overview

The `key_wrapper` package provides functionality to automatically distribute keys across multiple shards by appending postfixes like `:1`, `:2`, `:3`, etc. It supports dynamic shard count changes without requiring application restarts.

## Features

- **Automatic key distribution**: Evenly distributes keys across shards using cyclic postfix generation
- **Dynamic scaling**: Supports changing shard count at runtime
- **Two wrapper types**: General wrappers (update on any change) and growing-only wrappers (update only on increases)
- **Background monitoring**: Interrogator component for automatic shard count updates
- **Thread-safe**: All operations are safe for concurrent use

## Quick Start

```go
package main

import (
    "log"
    "time"
    "github.com/releaseband/wrappers/v2/key_wrapper"
)

func main() {
    // Create factory with 3 shards
    factory, err := key_wrapper.NewFactory(3)
    if err != nil {
        log.Fatal(err)
    }
    
    // Create a key wrapper
    wrapper := factory.MakeKeyWrapper()
    
    // Wrap keys - they'll get :1, :2, :3, :1, :2, :3, ...
    key1 := wrapper.WrapKey("user:123") // "user:123:1"
    key2 := wrapper.WrapKey("user:456") // "user:456:2"  
    key3 := wrapper.WrapKey("user:789") // "user:789:3"
    key4 := wrapper.WrapKey("user:999") // "user:999:1" (cycles back)
}
```

## Dynamic Shard Count Updates

Use the Interrogator to automatically update shard counts:

```go
// Function that returns current shard count
getShardsCount := func() (int, error) {
    // Your logic to determine current shard count
    // e.g., query database, check config, etc.
    return getCurrentShardCount()
}

// Configure interrogator
config := &key_wrapper.Config{
    GetShardsCount: getShardsCount,
    Factory:        factory,
    Interval:       30 * time.Second, // Check every 30 seconds
    ErrorHandler: func(err error) {
        log.Printf("Shard count update error: %v", err)
    },
}

// Start interrogator (runs in background)
interrogator, err := key_wrapper.RunInterrogator(config)
if err != nil {
    log.Fatal(err)
}

// Stop when done
defer interrogator.Stop()

// Or stop with context for timeout control
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
err = interrogator.StopWithContext(ctx)
if err != nil {
    log.Printf("Stop timeout: %v", err)
}
```

## Validation

The library includes comprehensive input validation:

- **Shard count**: Must be between 1 and 10,000
- **Configuration**: All required fields are validated
- **Runtime updates**: Invalid shard counts are rejected with descriptive errors

```go
// This will return an error
factory, err := key_wrapper.NewFactory(0) // Error: must be positive
factory, err := key_wrapper.NewFactory(20000) // Error: exceeds maximum
```

## Factory Statistics

Monitor your factory's state with built-in statistics:

```go
stats := factory.Stats()
fmt.Printf("Shards: %d\n", stats.Shards)
fmt.Printf("General Wrappers: %d\n", stats.GeneralWrappers)  
fmt.Printf("Growing-only Wrappers: %d\n", stats.GrowingWrappers)
```

## Wrapper Types

### General Wrapper
Updates on both increases and decreases in shard count:
```go
wrapper := factory.MakeKeyWrapper()
```

### Growing-Only Wrapper
Only updates when shard count increases (ignores decreases):
```go
wrapper := factory.MakeOnlyGrowingKeyWrapper()
```

## Use Cases

- **Redis Cluster**: Distribute keys across Redis cluster nodes
- **Database Sharding**: Distribute data across multiple database shards
- **Load Balancing**: Distribute requests across multiple service instances
- **Partitioned Systems**: Any system requiring even data distribution

## API Documentation

### KeyWrapper Interface
- `WrapKey(key string) string`: Wraps key with appropriate shard postfix

### WrapperFactory Interface
- `MakeKeyWrapper() KeyWrapper`: Creates general wrapper
- `MakeOnlyGrowingKeyWrapper() KeyWrapper`: Creates growing-only wrapper  
- `Stats() FactoryStats`: Returns factory statistics

### Factory
- `NewFactory(shardsCount int) (*Factory, error)`: Creates new factory with validation
- `MakeKeyWrapper() KeyWrapper`: Creates general wrapper
- `MakeOnlyGrowingKeyWrapper() KeyWrapper`: Creates growing-only wrapper
- `Stats() FactoryStats`: Returns current statistics

### FactoryStats
- `Shards int`: Current number of shards
- `GeneralWrappers int`: Number of general wrappers
- `GrowingWrappers int`: Number of growing-only wrappers

### Interrogator
- `RunInterrogator(cfg *Config) (*Interrogator, error)`: Starts background monitoring
- `Stop()`: Gracefully stops the interrogator
- `StopWithContext(ctx context.Context) error`: Stops with timeout control

### Config
- `GetShardsCount func() (int, error)`: Function to get current shard count
- `Factory *Factory`: Factory to update
- `Interval time.Duration`: Check interval
- `ErrorHandler func(err error)`: Required error handler

## Thread Safety

All components are designed for concurrent use:
- **Factory** uses RWMutex for safe shard count updates and wrapper management
- **KeyWrapper** uses Mutex for safe counter and shard count operations  
- **Interrogator** uses context-based cancellation and WaitGroup for clean shutdown
- **Store** is protected by Factory's mutex (no additional synchronization needed)
- **Interface compliance** is enforced at compile time with `var _ Interface = (*Implementation)(nil)` patterns

## Performance

The library is optimized for high-performance scenarios:
- Minimal memory allocations
- Lock-free reads where possible
- Efficient circular counter implementation
- Benchmark tests included for performance monitoring

## License

This project is part of the ReleaseBand ecosystem.