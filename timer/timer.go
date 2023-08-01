package timer

import (
	"time"
)

type timer struct {
	name     string
	duration time.Duration
	started  *time.Time
	elapsed  time.Duration
}

func New(options ...timerOption) timer {
	t := timer{}
	t.Config(options...)
	return t
}

func (t timer) Name() string {
	return t.name
}

func (t timer) Elapsed(now time.Time) time.Duration {
	e := t.elapsed
	if t.started != nil {
		e += now.Sub(*t.started)
	}
	return e
}

func (t *timer) Start(now time.Time) {
	if t.started == nil {
		t.started = &now
	}
}

func (t *timer) Stop(now time.Time) {
	if t.started != nil {
		t.elapsed += now.Sub(*t.started)
		t.started = nil
	}
}
