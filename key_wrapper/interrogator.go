// Package key_wrapper provides functionality for automatic key sharding with dynamic shard count adjustment.
package key_wrapper

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Interrogator periodically checks for shard count changes and updates the factory accordingly.
// It runs in the background and can be stopped when no longer needed.
type Interrogator struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// RunInterrogator starts a new interrogator with the given configuration.
// It runs the interrogator in a separate goroutine and returns a function to stop it.
// The returned stop function should be called when the interrogator is no longer needed
// to prevent goroutine leaks.
func RunInterrogator(cfg *Config) (*Interrogator, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validation config: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	srv := &Interrogator{
		cancel: cancel,
	}

	srv.wg.Add(1)

	go srv.run(ctx, cfg)

	return srv, nil
}

// Stop gracefully stops the interrogator and waits for it to finish.
// It's safe to call Stop multiple times.
func (l *Interrogator) Stop() {
	if l.cancel == nil {
		return // already stopped
	}

	l.cancel()
	l.wg.Wait()
}

func (l *Interrogator) StopWithContext(ctx context.Context) error {
	if l.cancel == nil {
		return nil // already stopped
	}

	done := make(chan struct{})

	go func() {
		l.cancel()
		l.wg.Wait()

		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func (l *Interrogator) run(ctx context.Context, cfg *Config) {
	defer l.wg.Done()

	t := time.NewTicker(cfg.Interval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			l.checkAndUpdate(cfg)
		}
	}
}

func (l *Interrogator) checkAndUpdate(cfg *Config) {
	count, err := cfg.GetShardsCount()
	if err != nil {
		cfg.ErrorHandler(err)
		return
	}

	err = cfg.Factory.compareAndUpdate(count)
	if err != nil {
		cfg.ErrorHandler(err)
		return
	}
}
