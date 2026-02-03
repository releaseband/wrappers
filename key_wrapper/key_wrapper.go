package key_wrapper

import (
	"strconv"
	"sync"
)

const (
	defaultPostfix = ":1"
)

// KeyWrapper provides functionality to wrap keys with shard postfixes.
// It automatically distributes keys across multiple shards by appending
// postfixes like ":1", ":2", etc.
type KeyWrapper interface {
	// WrapKey takes a base key and returns it with an appropriate shard postfix.
	// The postfix is determined by the current shard count and internal counter.
	WrapKey(key string) string
}

type WrapperFactory interface {
	// MakeKeyWrapper creates a new KeyWrapper that will be updated whenever
	// the factory's shard count changes (both increases and decreases).
	MakeKeyWrapper() KeyWrapper
	// MakeOnlyGrowingKeyWrapper creates a new KeyWrapper that will only
	// be updated when the factory's shard count increases.
	MakeOnlyGrowingKeyWrapper() KeyWrapper
	// Stats returns current statistics about the factory, including
	// the number of shards and registered wrappers.
	Stats() FactoryStats
}

// Compile-time interface compliance checks
var _ KeyWrapper = (*keyWrapper)(nil)
var _ WrapperFactory = (*Factory)(nil)

// keyWrapper is the concrete implementation of KeyWrapper interface.
// It maintains an internal counter (i) and current shard count to generate
// cyclic postfixes for even key distribution across shards.
type keyWrapper struct {
	mu          sync.Mutex // protects i and shardsCount from concurrent access
	i           int        // current position in the cycle (1 to shardsCount)
	shardsCount int        // total number of shards for distribution
}

// newKeyWrapper creates a new keyWrapper instance with the specified shard count.
// The wrapper starts with counter at 0 and will generate postfixes starting from ":1".
func newKeyWrapper(count int) *keyWrapper {
	w := &keyWrapper{}
	w.setCount(count)

	return w
}

// ResetShardsCount updates the shard count for this wrapper.
// This method is typically called by the factory when the global
// shard count changes. After calling this method, subsequent calls
// to WrapKey will use the new shard count for postfix generation.
func (b *keyWrapper) ResetShardsCount(count int) {
	b.setCount(count)
}

// setCount updates the shard count in a thread-safe manner.
// This method doesn't reset the counter position to maintain distribution consistency.
func (b *keyWrapper) setCount(count int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.shardsCount = count
}

// makePostfix generates the next shard postfix in the cycle.
// For single shard (shardsCount <= 1), it always returns ":1".
// For multiple shards, it increments the counter and wraps around when necessary.
// This method is thread-safe and ensures even distribution.
func (b *keyWrapper) makePostfix() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.shardsCount > 1 {
		b.i++
		if b.i > b.shardsCount {
			b.i = 1
		}
		return ":" + strconv.Itoa(b.i)
	}

	return defaultPostfix
}

// WrapKey wraps the given key with an appropriate shard postfix.
// For single shard (shardsCount <= 1), it always appends ":1".
// For multiple shards, it cycles through ":1", ":2", ..., ":shardsCount"
// to ensure even distribution across shards.
// Example: "user:123" -> "user:123:2"
func (b *keyWrapper) WrapKey(key string) string {
	return key + b.makePostfix()
}
