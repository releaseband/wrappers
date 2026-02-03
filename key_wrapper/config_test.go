package key_wrapper

import (
	"testing"
	"time"
)

func TestConfig_Validate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := &Config{
			GetShardsCount: func() (int, error) { return 1, nil },
			Factory:        &Factory{},
			Interval:       time.Second,
			ErrorHandler:   func(err error) {},
		}

		err := cfg.Validate()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("missing GetShardsCount", func(t *testing.T) {
		cfg := &Config{
			Factory:      &Factory{},
			Interval:     time.Second,
			ErrorHandler: func(err error) {},
		}

		err := cfg.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		expected := "GetShardsCount function is required"
		if err.Error() != expected {
			t.Fatalf("expected error %q, got %q", expected, err.Error())
		}
	})

	t.Run("missing Factory", func(t *testing.T) {
		cfg := &Config{
			GetShardsCount: func() (int, error) { return 1, nil },
			Interval:       time.Second,
			ErrorHandler:   func(err error) {},
		}

		err := cfg.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		expected := "Factory is required"
		if err.Error() != expected {
			t.Fatalf("expected error %q, got %q", expected, err.Error())
		}
	})

	t.Run("zero interval", func(t *testing.T) {
		cfg := &Config{
			GetShardsCount: func() (int, error) { return 1, nil },
			Factory:        &Factory{},
			Interval:       0,
			ErrorHandler:   func(err error) {},
		}

		err := cfg.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		expected := "Interval must be greater than zero"
		if err.Error() != expected {
			t.Fatalf("expected error %q, got %q", expected, err.Error())
		}
	})

	t.Run("negative interval", func(t *testing.T) {
		cfg := &Config{
			GetShardsCount: func() (int, error) { return 1, nil },
			Factory:        &Factory{},
			Interval:       -time.Second,
			ErrorHandler:   func(err error) {},
		}

		err := cfg.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		expected := "Interval must be greater than zero"
		if err.Error() != expected {
			t.Fatalf("expected error %q, got %q", expected, err.Error())
		}
	})

	t.Run("missing ErrorHandler", func(t *testing.T) {
		cfg := &Config{
			GetShardsCount: func() (int, error) { return 1, nil },
			Factory:        &Factory{},
			Interval:       time.Second,
		}

		err := cfg.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		expected := "ErrorHandler function is required"
		if err.Error() != expected {
			t.Fatalf("expected error %q, got %q", expected, err.Error())
		}
	})
}
