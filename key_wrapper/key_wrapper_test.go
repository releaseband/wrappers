package key_wrapper

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

func TestIncrementI(t *testing.T) {
	const shardsCount = 4

	kw := newKeyWrapper(shardsCount)
	exp := 0
	for i := 0; i < 100; i++ {
		kw.incrementI()
		exp++

		if exp > shardsCount {
			exp = 1
		}

		if kw.i != exp {
			t.Fatalf("exp:%d != got: %d", exp, kw.i)
		}
	}
}

func Test_makePostfix(t *testing.T) {
	const shardsCount = 5

	kw := newKeyWrapper(shardsCount)

	for i := 0; i < 5; i++ {

		j := i
		if j+1 > shardsCount {
			j = 0
		}

		exp := ":" + strconv.Itoa(j+1)
		got := kw.makePostfix()

		if exp != got {
			t.Fatal(fmt.Errorf("exp=%s, got=%s: %w", exp, got, errors.New("postfix invalid")))
		}
	}
}

func TestKeyWrapper_WrapKey(t *testing.T) {
	t.Run("shards count 0 or 1", func(t *testing.T) {

		check := func(shardsCount int) {
			kw := newKeyWrapper(shardsCount)

			key := "key"
			exp :=  key + ":1"
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


	{
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
}