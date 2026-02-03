package key_wrapper

// ResetShards defines the interface for objects that can have their
// shard count updated. This interface is used by the store to update
// all registered wrappers when shard count changes.
type ResetShards interface {
	// ResetShardsCount updates the shard count to the specified value.
	ResetShardsCount(count int)
}

type store struct {
	wrappers []ResetShards
}

func newStore() *store {
	return &store{
		wrappers: []ResetShards{},
	}
}

func (s *store) add(rs ResetShards) {
	s.wrappers = append(s.wrappers, rs)
}

func (s *store) update(count int) {
	for _, w := range s.wrappers {
		w.ResetShardsCount(count)
	}
}
