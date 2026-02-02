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
    "time"
    "github.com/releaseband/wrappers/key_wrapper"
)

func main() {
    // Create factory with 3 shards
    factory := key_wrapper.NewFactory(3)
    
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
}

// Start interrogator (runs in background)
stopFunc := key_wrapper.RunInterrogator(config)

// Stop when done
defer stopFunc()
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

### Factory
- `NewFactory(shardsCount int) *Factory`: Creates new factory
- `MakeKeyWrapper() KeyWrapper`: Creates general wrapper
- `MakeOnlyGrowingKeyWrapper() KeyWrapper`: Creates growing-only wrapper

### Interrogator
- `RunInterrogator(cfg *Config) func()`: Starts background monitoring
- `Stop()`: Stops the interrogator

## Thread Safety

All components are designed for concurrent use:
- Factory uses RWMutex for safe shard count updates
- Store uses Mutex for safe wrapper management
- Wrappers use atomic operations for counter management

## License

This project is part of the ReleaseBand ecosystem.