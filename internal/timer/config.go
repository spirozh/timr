package timer

import "time"

type TimerOption func(*Timer) TimerOption

func (t *Timer) Config(options ...TimerOption) []TimerOption {
	var initialOptions []TimerOption

	for _, option := range options {
		initialOptions = append(initialOptions, option(t))
	}

	return initialOptions
}

func Name(name string) TimerOption {
	return func(t *Timer) TimerOption {
		defer t.notify(func() { t.name = name })
		return Name(t.Name())
	}
}

func Duration(duration time.Duration) TimerOption {
	return func(t *Timer) TimerOption {
		defer t.notify(func() { t.duration = duration })
		return Duration(t.duration)
	}
}

func Started(started *time.Time) TimerOption {
	return func(t *Timer) TimerOption {
		defer t.notify(func() { t.started = started })
		return Started(t.started)
	}
}

func Elapsed(elapsed time.Duration) TimerOption {
	return func(t *Timer) TimerOption {
		defer t.notify(func() { t.elapsed = elapsed })
		return Elapsed(t.elapsed)
	}
}

func (t *Timer) notify(fn func()) {
	fn()
	// TODO: pub/sub
}
