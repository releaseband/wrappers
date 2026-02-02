package key_wrapper

import (
	"strconv"
)

const (
	defaultPostfix = ":1"
)

type keyWrapper struct {
	i           int
	shardsCount int
}

func (b *keyWrapper) setCount(count int) {
	b.shardsCount = count
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

func (b *keyWrapper) incrementI() int {
	b.i++

	if b.i > b.shardsCount {
		b.i = 1
	}

	return b.i
}

func (b *keyWrapper) makePostfix() string {
	if b.shardsCount > 1 {
		return ":" + strconv.Itoa(b.incrementI())
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
