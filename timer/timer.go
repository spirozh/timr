package timer

import (
	"time"

	"github.com/spirozh/timr"
)

type timer struct {
	clock  func() time.Time
	notify func(timr.TimrEventType, timr.Timer)

	duration time.Duration
	start    *time.Time    // when was timer resumed
	elapsed  time.Duration // how much time was elapsed before the timer was last started
}

var _ timr.Timer = (*timer)(nil)

func (t *timer) Resume() {
	if t.start == nil {
		now := t.clock()
		t.start = &now
	}

	if t.notify != nil {
		t.notify(timr.Resumed, t)
	}
}

func (t *timer) Pause() {
	if t.start != nil {
		t.elapsed += t.clock().Sub(*t.start)
		t.start = nil
	}

	if t.notify != nil {
		t.notify(timr.Paused, t)
	}
}

func (t *timer) Reset() {
	t.elapsed = 0
	if t.start != nil {
		*t.start = t.clock()
	}

	if t.notify != nil {
		t.notify(timr.Reset, t)
	}
}

func (t *timer) State() timr.TimerState {
	ts := timr.TimerState{
		Duration:  t.duration.Milliseconds(),
		Remaining: (t.duration - t.elapsed).Milliseconds(),
	}

	if t.start != nil {
		ts.Running = true
		ts.Remaining -= t.clock().Sub(*t.start).Milliseconds()
	}

	return ts
}
