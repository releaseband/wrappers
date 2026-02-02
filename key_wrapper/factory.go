package key_wrapper

import (
	"fmt"
	"sync"
)

// Factory creates and manages KeyWrapper instances.
// It maintains two types of wrappers: general wrappers that update on any shard count change,
// and only-growing wrappers that only update when shard count increases.
// Factory ensures thread-safe operations and shard count management.
type Factory struct {
	mu                  *sync.RWMutex
	generalWrappers     *store
	onlyGrowingWrappers *store
	shardsCount         int
}

type FactoryStats struct {
	Shards          int
	GeneralWrappers int
	GrowingWrappers int
}

const (
	minShardsCount = 1
	maxShardsCount = 10_000
)

// NewFactory creates a new Factory with the specified initial shard count.
// The shard count determines how many different postfixes will be used
// when wrapping keys (e.g., ":1", ":2", ":3" for shardsCount=3).
// An error is returned if the initial shard count is
// less than 1 and greater than 10_000
func NewFactory(initialShardsCount int) (*Factory, error) {
	if err := validateShardsCount(initialShardsCount); err != nil {
		return nil, err
	}

	return &Factory{
		mu:                  &sync.RWMutex{},
		onlyGrowingWrappers: newStore(),
		generalWrappers:     newStore(),
		shardsCount:         initialShardsCount,
	}, nil
}

func validateShardsCount(count int) error {
	if count < minShardsCount {
		return fmt.Errorf("initial shards count must be positive, got %d",
			count)
	}

	if count > maxShardsCount {
		return fmt.Errorf("initial shards count must be less than %d, got %d",
			maxShardsCount, count)
	}

	return nil
}

// MakeKeyWrapper creates a new KeyWrapper that will be updated whenever
// the factory's shard count changes (both increases and decreases).
// The returned wrapper is registered with the factory and will automatically
// receive shard count updates.
func (f *Factory) MakeKeyWrapper() KeyWrapper {
	f.mu.Lock()
	defer f.mu.Unlock()

	w := newKeyWrapper(f.shardsCount)
	f.generalWrappers.add(w)

	return w
}

// MakeOnlyGrowingKeyWrapper creates a new KeyWrapper that will only be updated
// when the factory's shard count increases, but not when it decreases.
// This is useful for scenarios where reducing shard count should not affect
// existing key distribution patterns.
func (f *Factory) MakeOnlyGrowingKeyWrapper() KeyWrapper {
	f.mu.Lock()
	defer f.mu.Unlock()

	w := newKeyWrapper(f.shardsCount)
	f.onlyGrowingWrappers.add(w)

	return w
}

func (f *Factory) compareAndUpdate(shardCount int) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if shardCount == f.shardsCount {
		// No change in shard count
		return nil
	}

	err := validateShardsCount(shardCount)
	if err != nil {
		return err
	}

	f.generalWrappers.update(shardCount)

	if shardCount > f.shardsCount {
		f.onlyGrowingWrappers.update(shardCount)
	}

	f.shardsCount = shardCount

	return nil
}

func (f *Factory) Stats() FactoryStats {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return FactoryStats{
		Shards:          f.shardsCount,
		GeneralWrappers: len(f.generalWrappers.wrappers),
		GrowingWrappers: len(f.onlyGrowingWrappers.wrappers),
	}
}
