package key_wrapper

import "time"

type Config struct {
	ErrorHandler   func(mess string, err error)
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
		if err != nil {
			cfg.ErrorHandler("get shards count failed", err)
		} else {
			cfg.Factory.updateShardsCount(count)
		}
	}
}

func RunInterrogator(cfg *Config) func() {
	l :=  &Interrogator{}
	l.run(cfg)

	return l.Stop
}

func (l *Interrogator) Stop() {
	if l.stopTick != nil {
		l.stopTick()
	}
}