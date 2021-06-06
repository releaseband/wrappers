package key_wrapper

import "sync"

type KeyWrapper interface {
	WrapKey(key string) string
}

type Factory struct {
	generalWrappers     *store
	onlyGrowingWrappers *store
	mu                  *sync.RWMutex
	shardsCount         int
}

func NewFactory(slotsCount int) *Factory {
	return &Factory{
		mu:          &sync.RWMutex{},
		shardsCount: slotsCount,
	}
}

func (f *Factory) makeKeyWrapper() *keyWrapper {
	return newKeyWrapper(f.shardsCount)
}

func (f *Factory) MakeKeyWrapper() KeyWrapper {
	w := f.makeKeyWrapper()

	if f.generalWrappers == nil {
		f.generalWrappers = newStore()
	}

	f.generalWrappers.add(w)

	return w
}

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
	var decreased bool

	f.mu.Lock()
	if shardCount < f.shardsCount {
		decreased = true
	}

	f.shardsCount = shardCount
	f.mu.Unlock()

	f.generalWrappers.update(shardCount)

	if !decreased {
		f.onlyGrowingWrappers.update(shardCount)
	}
}

func (f *Factory) UpdateShardsCount(shardsCount int) {
	if f.getCount() != shardsCount {
		f.updateShardsCount(shardsCount)
	}
}
