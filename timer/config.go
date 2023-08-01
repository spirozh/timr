package timer

import "time"

type timerOption func(*timer) timerOption

func (t *timer) Config(options ...timerOption) []timerOption {
	var starting []timerOption

	for _, option := range options {
		starting = append(starting, option(t))
	}

	return starting
}

func Name(name string) timerOption {
	return func(t *timer) timerOption {
		defer func() { t.name = name }()
		return Name(t.Name())
	}
}

func Duration(duration time.Duration) timerOption {
	return func(t *timer) timerOption {
		defer func() { t.duration = duration }()
		return Duration(t.duration)
	}
}

func Started(started *time.Time) timerOption {
	return func(t *timer) timerOption {
		defer func() { t.started = started }()
		return Started(t.started)
	}
}

func Elapsed(elapsed time.Duration) timerOption {
	return func(t *timer) timerOption {
		defer func() { t.elapsed = elapsed }()
		return Elapsed(t.elapsed)
	}
}
