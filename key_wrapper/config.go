package key_wrapper

import (
	"errors"
	"time"
)

// Config holds the configuration for the Interrogator.
// It defines how the interrogator should check for shard
// count changes and update the factory.
type Config struct {
	// GetShardsCount is a function that returns the current number of shards.
	// It should return an error if the shard count cannot be determined.
	GetShardsCount func() (int, error)
	// Factory is the factory instance that will be updated with new shard counts.
	Factory *Factory
	// Interval specifies how often the interrogator
	// should check for shard count changes.
	Interval time.Duration

	// ErrorHandler is a required function used to handle errors
	// encountered during shard count retrieval.
	ErrorHandler func(err error)
}

func (cfg *Config) Validate() error {
	if cfg.GetShardsCount == nil {
		return errors.New("GetShardsCount function is required")
	}

	if cfg.Factory == nil {
		return errors.New("Factory is required")
	}

	if cfg.Interval <= 0 {
		return errors.New("Interval must be greater than zero")
	}

	if cfg.ErrorHandler == nil {
		return errors.New("ErrorHandler function is required")
	}

	return nil
}
