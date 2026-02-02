package key_wrapper

import (
	"strconv"
	"sync"
	"testing"
)

func TestFactory(t *testing.T) {
	const shardsCount = 4

	f, err := NewFactory(shardsCount)
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}

	if f.shardsCount != shardsCount {
		t.Fatal("shardsCount invalid")
	}

	if f.shardsCount != shardsCount {
		t.Fatal("getCount invalid")
	}

	if f.onlyGrowingWrappers == nil {
		t.Fatal("onlyGrowingWrappers should not be nil")
	}

	if f.generalWrappers == nil {
		t.Fatal("generalWrappers should not be nil")
	}

	generalWrappers := make([]KeyWrapper, 6)
	for i := 0; i < 6; i++ {
		generalWrappers[i] = f.MakeKeyWrapper()
	}

	if f.generalWrappers == nil {
		t.Fatal("generalWrappers should not be nil")
	}

	if len(f.generalWrappers.wrappers) != 6 {
		t.Fatal("generalWrappers.wrappers len should be 6")
	}

	onlyGrowingWrappers := make([]KeyWrapper, 3)
	for i := 0; i < 3; i++ {
		onlyGrowingWrappers[i] = f.MakeOnlyGrowingKeyWrapper()
	}

	if len(f.generalWrappers.wrappers) != 6 {
		t.Fatal("generalWrappers.wrappers len should be 6")
	}

	if len(f.onlyGrowingWrappers.wrappers) != 3 {
		t.Fatal("onlyGrowingWrappers.wrappers len should be 3")
	}

	f.compareAndUpdate(shardsCount - 1)

	const key = "KEY"
	expKey := key + ":1"

	for _, generalW := range generalWrappers {
		for i := 0; i < shardsCount-1; i++ {
			generalW.WrapKey(key)
		}

		gotKey := generalW.WrapKey(key)
		if expKey != gotKey {
			t.Fatalf("gotKey=%s, expKey=%s", gotKey, expKey)
		}
	}

	expKey = key + ":" + strconv.Itoa(shardsCount)
	for _, onlyGrowing := range onlyGrowingWrappers {
		for i := 0; i < shardsCount-1; i++ {
			onlyGrowing.WrapKey(key)
		}

		gotKey := onlyGrowing.WrapKey(key)
		if gotKey != expKey {
			t.Fatalf("gotKey(%s) != expKey(%s)", gotKey, expKey)
		}
	}
}

func BenchmarkWrapKey(b *testing.B) {
	factory, _ := NewFactory(10)
	wrapper := factory.MakeKeyWrapper()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wrapper.WrapKey("testkey")
	}
}

func TestConcurrentAccess(t *testing.T) {
	factory, _ := NewFactory(5)
	wrapper := factory.MakeKeyWrapper()

	count := 100
	wg := sync.WaitGroup{}
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			wrapper.WrapKey("test")
		}()
	}
	wg.Wait()
}
