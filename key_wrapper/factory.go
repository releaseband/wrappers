package key_wrapper

import (
	"fmt"
	"sync"
)

// Factory creates and manages KeyWrapper instances.
// It maintains two types of wrappers: general wrappers that update on any shard count change,
// and only-growing wrappers that only update when shard count increases.
// Factory ensures thread-safe operations and shard count management.
//
// All public methods are thread-safe and can be called concurrently.
type Factory struct {
	mu                  *sync.RWMutex // protects all fields from concurrent access
	generalWrappers     *store        // wrappers that update on any shard count change
	onlyGrowingWrappers *store        // wrappers that only update on shard count increases
	shardsCount         int           // current number of shards for key distribution
}

// FactoryStats provides statistical information about a Factory instance.
// All values represent the current state at the time of the Stats() call.
type FactoryStats struct {
	Shards          int // current number of shards configured
	GeneralWrappers int // number of registered general wrappers
	GrowingWrappers int // number of registered growing-only wrappers
}

const (
	// minShardsCount defines the minimum allowed number of shards.
	// At least one shard is required for key distribution.
	minShardsCount = 0
	// maxShardsCount defines the maximum allowed number of shards.
	// This limit prevents excessive memory usage and ensures reasonable performance.
	maxShardsCount = 10_000
)

// NewFactory creates a new Factory with the specified initial shard count.
// The shard count determines how many different postfixes will be used
// when wrapping keys (e.g., ":1", ":2", ":3" for shardsCount=3).
// An error is returned if the initial shard count is
// less than 0 or greater than 10_000
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

// validateShardsCount checks if the provided shard count is within acceptable limits.
// Returns an error with a descriptive message if the count is invalid.
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

// compareAndUpdate updates the factory's shard count if it differs from the new value.
// This method is called by the Interrogator to apply shard count changes.
// It validates the new count and updates appropriate wrappers based on their type:
// - General wrappers are always updated
// - Growing-only wrappers are updated only when shard count increases
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

// Stats returns current statistics about the factory.
// The returned FactoryStats contains information about shard count
// and the number of registered wrappers of each type.
// This method is thread-safe and provides a consistent snapshot.
func (f *Factory) Stats() FactoryStats {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return FactoryStats{
		Shards:          f.shardsCount,
		GeneralWrappers: len(f.generalWrappers.wrappers),
		GrowingWrappers: len(f.onlyGrowingWrappers.wrappers),
	}
}
