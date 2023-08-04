package timer

import (
	"fmt"
	"time"
)

type Timer struct {
	name     string
	duration time.Duration
	started  *time.Time
	elapsed  time.Duration
}

func New(options ...timerOption) Timer {
	t := Timer{}
	t.Config(options...)
	return t
}

func (t Timer) String() string {
	s := "nil"
	if t.started != nil {
		s = fmt.Sprintf(`"%s"`, t.started.Format(timeFormat))
	}
	return fmt.Sprintf(`<timer "%s" duration:%s started:%s elapsed:%s>`, t.name, t.duration, s, t.elapsed)
}

func (t Timer) Name() string {
	return t.name
}

func (t Timer) Elapsed(now time.Time) time.Duration {
	e := t.elapsed
	if t.started != nil {
		e += now.Sub(*t.started)
	}
	return e
}

func (t Timer) Remaining(now time.Time) time.Duration {
	return t.duration - t.Elapsed(now)
}

func (t *Timer) Start(now time.Time) {
	if t.started == nil {
		t.started = &now
	}
}

func (t *Timer) Stop(now time.Time) {
	if t.started != nil {
		t.elapsed += now.Sub(*t.started)
		t.started = nil
	}
}
