package key_wrapper

import (
	"strconv"
	"testing"
)

func TestKeyPostfix_Next(t *testing.T) {
	const key = "key"

	t.Run("shards count <= 1", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			kp := NewKeyWrapper(i)
			exp := key + defaultPostfix
			for i := 0; i < 100; i++ {
				got := kp.WrapKey(key)

				if got != exp {
					t.Fatalf("exp=%s | got=%s", exp, got)
				}
			}
		}
	})

	t.Run("count > 1", func(t *testing.T) {
		const count = 10
		kp := NewKeyWrapper(count)

		if kp.shardsCount != count {
			t.Fatalf("count should be equal %d", count)
		}

		if len(kp.postfixes) != count {
			t.Fatalf("postfixes count should be equeal %d", count)
		}

		index := 1
		for i := 0; i < count*2-1; i++ {
			index++

			if i == count-1 {
				index =1
			}

			exp := key + ":" + strconv.Itoa(index)
			got := kp.WrapKey(key)
			if got != exp {
				t.Fatalf("exp=%s | got=%s", exp, got)
			}
		}
	})
}
