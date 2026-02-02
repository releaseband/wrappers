// Package key_wrapper provides functionality for automatic key sharding with dynamic shard count adjustment.
package key_wrapper

import "time"

// Config holds the configuration for the Interrogator.
// It defines how the interrogator should check for shard count changes and update the factory.
type Config struct {
	// GetShardsCount is a function that returns the current number of shards.
	// It should return an error if the shard count cannot be determined.
	GetShardsCount func() (int, error)
	// Factory is the factory instance that will be updated with new shard counts.
	Factory        *Factory
	// Interval specifies how often the interrogator should check for shard count changes.
	Interval       time.Duration
}

// Interrogator periodically checks for shard count changes and updates the factory accordingly.
// It runs in the background and can be stopped when no longer needed.
type Interrogator struct {
	stopTick func()
}

func (l *Interrogator) run(cfg *Config) {
	t := time.NewTicker(cfg.Interval)

	l.stopTick = t.Stop

	for range t.C {
		count, err := cfg.GetShardsCount()
		if err == nil {
			cfg.Factory.updateShardsCount(count)
		}
	}
}

// RunInterrogator starts a new interrogator with the given configuration.
// It runs the interrogator in a separate goroutine and returns a function to stop it.
// The returned stop function should be called when the interrogator is no longer needed
// to prevent goroutine leaks.
func RunInterrogator(cfg *Config) func() {
	l := &Interrogator{}
	go l.run(cfg)

	return l.Stop
}

// Stop gracefully stops the interrogator by stopping its ticker.
// It's safe to call Stop multiple times.
func (l *Interrogator) Stop() {
	if l.stopTick != nil {
		l.stopTick()
	}
}
