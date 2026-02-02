package key_wrapper_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/releaseband/wrappers/key_wrapper"
)

// Test_Example demonstrates basic usage of the key_wrapper library
func TestExample(t *testing.T) {
	// Create factory with 3 shards
	factory, err := key_wrapper.NewFactory(3)
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}

	// Create a key wrapper
	wrapper := factory.MakeKeyWrapper()

	// Wrap several keys - they'll cycle through :1, :2, :3
	outputs := []string{
		"user:100:1",
		"user:200:2",
		"user:300:3",
		"user:400:1",
		"user:500:2",
		"user:600:3",
	}

	for i := 1; i <= 6; i++ {
		key := fmt.Sprintf("user:%d", i*100)
		wrappedKey := wrapper.WrapKey(key)
		expected := outputs[i-1]
		if wrappedKey != expected {
			t.Errorf("Expected %s, got %s", expected, wrappedKey)
		}
	}
}

// Test_ExampleInterrogator demonstrates how to use the interrogator for dynamic shard count updates
func Test_ExampleKeyWrapper(t *testing.T) {
	// Create factory
	const (
		initialShards = 2
		updatedShards = 4
	)

	factory, err := key_wrapper.NewFactory(initialShards)
	if err != nil {
		t.Fatalf("failed to create factory: %v", err)
	}

	// Mock function that returns current shard count
	// In real usage, this would query your infrastructure
	var callCount int
	getCurrentShards := func() (int, error) {
		callCount++
		// First few calls return 2 shards, later calls return 4 shards
		if callCount <= 2 {
			return initialShards, nil
		}
		return updatedShards, nil
	}

	// Configure interrogator with short interval for demo
	config := &key_wrapper.Config{
		GetShardsCount: getCurrentShards,
		Factory:        factory,
		Interval:       100 * time.Millisecond, // Check every 100ms for demo
		ErrorHandler: func(err error) {
			t.Fatal("error in interrogator: ", err)
		},
	}

	// Start interrogator in background
	srv, err := key_wrapper.RunInterrogator(config)
	if err != nil {
		t.Fatalf("should be not error: %v", err)
	}

	defer srv.Stop() // Important: always stop the interrogator when done

	// Create wrapper
	wrapper := factory.MakeKeyWrapper()

	outputs := []string{
		"key1:1",
		"key2:2",
		"key3:1",
		"key4:2",
		"key5:3",
		"key6:4",
		"key7:1",
	}

	t.Log("Before interrogator detects change:")
	for i, exp := range outputs {
		if i == 3 {
			time.Sleep(300 * time.Millisecond)
			t.Log("After interrogator detects change to 4 shards:")
		}

		got := wrapper.WrapKey("key" + strconv.Itoa(i+1))
		if got != exp {
			t.Fatalf("exp=%s, got=%s", exp, got)
		}
	}
}

// Test_ExampleOnlyGrowKeyWrapper demonstrates usage of only-growing key wrappers
func Test_ExampleOnlyGrowKeyWrapper(t *testing.T) {
	// Create factory
	const (
		initialShards = 4
		updatedShards = 2
	)

	factory, err := key_wrapper.NewFactory(initialShards)
	if err != nil {
		t.Fatalf("should be not error: %v", err)
	}

	// Mock function that returns current shard count
	// In real usage, this would query your infrastructure
	var callCount int
	getCurrentShards := func() (int, error) {
		callCount++
		// First few calls return 4 shards,
		if callCount <= 2 {
			return initialShards, nil
		}
		// later calls return 2 shards
		return updatedShards, nil
	}

	// Configure interrogator with short interval for demo
	config := &key_wrapper.Config{
		GetShardsCount: getCurrentShards,
		Factory:        factory,
		Interval:       100 * time.Millisecond, // Check every 100ms for demo
		ErrorHandler: func(err error) {
			t.Fatal("error in interrogator:", err)
		},
	}

	// Start interrogator in background
	srv, err := key_wrapper.RunInterrogator(config)
	if err != nil {
		t.Fatalf("should be not error: %v", err)
	}

	defer srv.Stop() // Important: always stop the interrogator when done

	// Create wrapper
	wrapper := factory.MakeOnlyGrowingKeyWrapper()

	// Expected outputs - should continue using 4 shards even after decrease
	// in shard count
	// because it's an only-growing wrapper
	outputs := []string{
		"key1:1",
		"key2:2",
		"key3:3",
		"key4:4",
		"key5:1",
		"key6:2",
		"key7:3",
		"key8:4",
	}

	t.Log("Before interrogator detects change:")
	for i, exp := range outputs {
		if i == 3 {
			time.Sleep(300 * time.Millisecond)
			t.Log("After interrogator detects change to 4 shards:")
		}

		got := wrapper.WrapKey("key" + strconv.Itoa(i+1))
		if got != exp {
			t.Fatalf("exp=%s, got=%s", exp, got)
		}
	}
}
