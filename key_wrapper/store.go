package key_wrapper

import "sync"

// ResetShards defines the interface for objects that can have their
// shard count updated. This interface is used by the store to update
// all registered wrappers when shard count changes.
type ResetShards interface {
	// ResetShardsCount updates the shard count to the specified value.
	ResetShardsCount(count int)
}

type store struct {
	mu       *sync.Mutex
	wrappers []ResetShards
}

func newStore() *store {
	return &store{
		mu:       &sync.Mutex{},
		wrappers: []ResetShards{},
	}
}

func (s *store) add(rs ResetShards) {
	s.mu.Lock()
	s.wrappers = append(s.wrappers, rs)
	s.mu.Unlock()
}

func (s *store) update(count int) {
	s.mu.Lock()

	for _, w := range s.wrappers {
		w.ResetShardsCount(count)
	}

	s.mu.Unlock()
}
