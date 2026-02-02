package key_wrapper

import (
	"sync"
)

// KeyWrapper provides functionality to wrap keys with shard postfixes.
// It automatically distributes keys across multiple shards by appending
// postfixes like ":1", ":2", etc.
type KeyWrapper interface {
	// WrapKey takes a base key and returns it with an appropriate shard postfix.
	// The postfix is determined by the current shard count and internal counter.
	WrapKey(key string) string
}

// Factory creates and manages KeyWrapper instances.
// It maintains two types of wrappers: general wrappers that update on any shard count change,
// and only-growing wrappers that only update when shard count increases.
type Factory struct {
	generalWrappers     *store
	onlyGrowingWrappers *store
	mu                  *sync.RWMutex
	shardsCount         int
}

// NewFactory creates a new Factory with the specified initial shard count.
// The shard count determines how many different postfixes will be used
// when wrapping keys (e.g., ":1", ":2", ":3" for shardsCount=3).
func NewFactory(shardsCount int) *Factory {
	return &Factory{
		mu:          &sync.RWMutex{},
		shardsCount: shardsCount,
		//initialized lazily to avoid unnecessary allocations
		onlyGrowingWrappers: nil,
		//initialized lazily to avoid unnecessary allocations
		generalWrappers: nil,
	}
}

func (f *Factory) makeKeyWrapper() *keyWrapper {
	return newKeyWrapper(f.shardsCount)
}

// MakeKeyWrapper creates a new KeyWrapper that will be updated whenever
// the factory's shard count changes (both increases and decreases).
// The returned wrapper is registered with the factory and will automatically
// receive shard count updates.
func (f *Factory) MakeKeyWrapper() KeyWrapper {
	w := f.makeKeyWrapper()

	if f.generalWrappers == nil {
		f.generalWrappers = newStore()
	}

	f.generalWrappers.add(w)

	return w
}

// MakeOnlyGrowingKeyWrapper creates a new KeyWrapper that will only be updated
// when the factory's shard count increases, but not when it decreases.
// This is useful for scenarios where reducing shard count should not affect
// existing key distribution patterns.
func (f *Factory) MakeOnlyGrowingKeyWrapper() KeyWrapper {
	w := f.makeKeyWrapper()

	if f.onlyGrowingWrappers == nil {
		f.onlyGrowingWrappers = newStore()
	}

	f.onlyGrowingWrappers.add(w)

	return w
}

func (f *Factory) getCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.shardsCount
}

func (f *Factory) updateShardsCount(shardCount int) {
	if shardCount == f.getCount() {
		return
	}

	var decreased bool

	f.mu.Lock()
	if shardCount < f.shardsCount {
		decreased = true
	}

	f.shardsCount = shardCount
	f.mu.Unlock()

	if f.generalWrappers != nil {
		f.generalWrappers.update(shardCount)
	}

	if !decreased && f.onlyGrowingWrappers != nil {
		f.onlyGrowingWrappers.update(shardCount)
	}
}
