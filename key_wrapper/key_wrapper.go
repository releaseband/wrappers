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

func (b *keyWrapper) WrapKey(key string) string {
	return key + b.makePostfix()
}
