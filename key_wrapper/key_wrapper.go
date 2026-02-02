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
	// GetStats returns current statistics about the factory, including
	// the number of shards and registered wrappers.
	Stats() FactoryStats
}

var _ KeyWrapper = (*keyWrapper)(nil)
var _ WrapperFactory = (*Factory)(nil)

type keyWrapper struct {
	mu          sync.Mutex
	i           int
	shardsCount int
}

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

func (b *keyWrapper) setCount(count int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.shardsCount = count
}

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
