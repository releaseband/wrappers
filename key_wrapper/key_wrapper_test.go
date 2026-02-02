package key_wrapper

import (
	"strconv"
	"testing"
)

func TestWrapKey(t *testing.T) {
	const shardsCount = 4

	kw := newKeyWrapper(shardsCount)

	// Test cycling through shard postfixes
	expectedPostfixes := []string{":1", ":2", ":3", ":4", ":1", ":2", ":3", ":4"}

	for i, expectedPostfix := range expectedPostfixes {
		result := kw.WrapKey("test")
		expected := "test" + expectedPostfix

		if result != expected {
			t.Fatalf("iteration %d: expected %s, got %s", i, expected, result)
		}
	}
}

func TestKeyWrapper_WrapKey(t *testing.T) {
	t.Run("shards count 0 or 1", func(t *testing.T) {

		check := func(shardsCount int) {
			kw := newKeyWrapper(shardsCount)

			key := "key"
			exp := key + ":1"
			for i := 0; i < 100; i++ {
				wrappedKey := kw.WrapKey(key)
				if exp != wrappedKey {
					t.Fatalf("got=%s != exp=%s", wrappedKey, exp)
				}
			}
		}

		check(0)
		check(1)
	})

	t.Run("shards count > 1", func(t *testing.T) {
		const shardsCount = 6

		kw := newKeyWrapper(shardsCount)
		key := "key"

		var j int
		for i := 0; i < 7; i++ {
			j++

			if j > shardsCount {
				j = 1
			}

			expWrappedKey := key + ":" + strconv.Itoa(j)
			gotWrappedKey := kw.WrapKey(key)

			if gotWrappedKey != expWrappedKey {
				t.Fatalf("gotKey(%s) != expKey(%s)", gotWrappedKey, expWrappedKey)
			}
		}
	})
}
