package timer

import (
	"fmt"
	"time"
)

type Timer interface {
	Config(options ...timerOption) []timerOption
	Name() string
	Elapsed(time.Time) time.Duration
	Remaining(time.Time) time.Duration
	Start(time.Time)
	Stop(time.Time)
}

type timer struct {
	name     string
	duration time.Duration
	started  *time.Time
	elapsed  time.Duration
}

func New(options ...timerOption) Timer {
	t := timer{}
	t.Config(options...)
	return &t
}

func (t timer) String() string {
	s := "nil"
	if t.started != nil {
		s = fmt.Sprintf(`"%s"`, t.started.Format(timeFormat))
	}
	return fmt.Sprintf(`<timer "%s" duration:%s started:%s elapsed:%s>`, t.name, t.duration, s, t.elapsed)
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

func (t timer) Remaining(now time.Time) time.Duration {
	return t.duration - t.Elapsed(now)
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
