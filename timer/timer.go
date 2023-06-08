package timer

import (
	"time"

	"github.com/spirozh/timr"
)

type timer struct {
	now    func() time.Time
	notify func(timr.TimrEventType, timr.Timer)

	duration time.Duration
	start    *time.Time
	elapsed  time.Duration
}

var _ timr.Timer = (*timer)(nil)

func (t *timer) Resume() {
	if t.start == nil {
		now := t.now()
		t.start = &now
	}

	t.notify(timr.EventTimerResumed, t)
}

func (t *timer) Pause() {
	if t.start != nil {
		t.elapsed += t.now().Sub(*t.start)
		t.start = nil
	}

	t.notify(timr.EventTimerPaused, t)
}

func (t *timer) Reset() {
	t.start, t.elapsed = nil, 0

	t.notify(timr.EventTimerReset, t)
}

func (t *timer) Remaining() (remaining time.Duration, isRunning bool) {
	remaining = t.duration - t.elapsed

	if t.start == nil {
		return remaining, false
	}

	return remaining - t.now().Sub(*t.start), true
}
