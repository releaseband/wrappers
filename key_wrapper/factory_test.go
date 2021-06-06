package key_wrapper

import (
	"strconv"
	"testing"
)

func TestFactory(t *testing.T) {
	const shardsCount = 4
	f := NewFactory(shardsCount)
	if f.shardsCount != shardsCount {
		t.Fatal("shardsCount invalid")
	}

	count := f.getCount()
	if count != shardsCount {
		t.Fatal("getCount invalid")
	}

	if f.onlyGrowingWrappers != nil {
		t.Fatal("onlyGrowingWrappers should be nil")
	}

	if f.generalWrappers != nil {
		t.Fatal("generalWrappers should be nil")
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

	if f.onlyGrowingWrappers != nil {
		t.Fatal("onlyGrowingWrappers should be nil")
	}

	onlyGrowingWrappers := make([]KeyWrapper, 3)
	for i := 0; i < 3; i++ {
		onlyGrowingWrappers[i] = f.MakeOnlyGrowingKeyWrapper()
	}

	if len(f.generalWrappers.wrappers) != 6 {
		t.Fatal("generalWrappers.wrappers len should be 6")
	}

	if f.onlyGrowingWrappers == nil {
		t.Fatal("onlyGrowingWrappers should not be nil")
	}

	if len(f.onlyGrowingWrappers.wrappers) != 3 {
		t.Fatal("onlyGrowingWrappers.wrappers len should be 3")
	}

	f.UpdateShardsCount(shardsCount - 1)
	const key = "KEY"
	expKey := key + ":1"

	for _, generalW := range generalWrappers {
		for i := 0; i < shardsCount - 1; i++ {
			generalW.WrapKey(key)
		}

		gotKey := generalW.WrapKey(key)
		if expKey != gotKey {
			t.Fatalf("gotKey=%s, expKey=%s", gotKey, expKey)
		}
	}

	expKey = key + ":" + strconv.Itoa(shardsCount)
	for _, onlyGrowing := range onlyGrowingWrappers {
		for i := 0; i < shardsCount - 1; i++ {
			onlyGrowing.WrapKey(key)
		}

		gotKey := onlyGrowing.WrapKey(key)
		if gotKey != expKey {
			t.Fatalf("gotKey(%s) != expKey(%s)", gotKey, expKey)
		}
	}
}
