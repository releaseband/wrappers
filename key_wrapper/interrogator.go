package key_wrapper

import "time"

type Config struct {
	GetShardsCount func() (int, error)
	Factory        *Factory
	Interval       time.Duration
}

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

func RunInterrogator(cfg *Config) func() {
	l := &Interrogator{}
	go l.run(cfg)

	return l.Stop
}

func (l *Interrogator) Stop() {
	if l.stopTick != nil {
		l.stopTick()
	}
}
