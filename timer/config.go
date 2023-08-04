package timer

import "time"

type timerOption func(*Timer) timerOption

func (t *Timer) Config(options ...timerOption) []timerOption {
	var initialOptions []timerOption

	for _, option := range options {
		initialOptions = append(initialOptions, option(t))
	}

	return initialOptions
}

func Name(name string) timerOption {
	return func(t *Timer) timerOption {
		defer t.notify(func() { t.name = name })
		return Name(t.Name())
	}
}

func Duration(duration time.Duration) timerOption {
	return func(t *Timer) timerOption {
		defer t.notify(func() { t.duration = duration })
		return Duration(t.duration)
	}
}

func Started(started *time.Time) timerOption {
	return func(t *Timer) timerOption {
		defer t.notify(func() { t.started = started })
		return Started(t.started)
	}
}

func Elapsed(elapsed time.Duration) timerOption {
	return func(t *Timer) timerOption {
		defer t.notify(func() { t.elapsed = elapsed })
		return Elapsed(t.elapsed)
	}
}

func (t *Timer) notify(fn func()) {
	fn()
	// TODO: pub/sub
}
