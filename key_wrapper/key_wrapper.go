package key_wrapper

import (
	"strconv"
)

const (
	defaultPostfix = ":1"
)

type KeyWrapper struct {
	i           int
	shardsCount int
	postfixes   []string
}

func getPostfix(i int) string {
	return ":" + strconv.Itoa(i)
}

func makePostfixes(count int) []string {
	postfixes := make([]string, count)
	for i := 0; i < count; i++ {
		postfixes[i] = getPostfix(i+1)
	}

	return postfixes
}

func (b *KeyWrapper) setCount(count int) {
	b.shardsCount = count
	b.postfixes = makePostfixes(count)
}

func NewKeyWrapper(count int) *KeyWrapper {
	w := &KeyWrapper{}
	w.setCount(count)

	return w
}

func (b *KeyWrapper) ResetShardsCount(count int) {
	b.setCount(count)
}

func (b *KeyWrapper) WrapKey(key string) string {
	var postfix string
	if b.shardsCount > 1 {
		b.i++
		if b.i >= b.shardsCount {
			b.i = 0
		}

		postfix = b.postfixes[b.i]
	} else {
		postfix = defaultPostfix
	}

	return key + postfix
}
