package key_wrapper

import "sync"

type ResetShards interface {
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
