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
		t.notify(timr.EventTimerResumed, t)
	}
}

func (t *timer) Pause() {
	if t.start != nil {
		t.elapsed += t.clock().Sub(*t.start)
		t.start = nil
	}

	if t.notify != nil {
		t.notify(timr.EventTimerPaused, t)
	}
}

func (t *timer) Reset() {
	t.elapsed = 0
	if t.start != nil {
		*t.start = t.clock()
	}

	if t.notify != nil {
		t.notify(timr.EventTimerReset, t)
	}
}

func (t *timer) Remaining() (remaining time.Duration, isRunning bool) {
	remaining = t.duration - t.elapsed

	if t.start == nil {
		return remaining, false
	}

	return remaining - t.clock().Sub(*t.start), true
}
