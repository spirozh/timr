package timer

import (
	"time"

	"github.com/spirozh/timr"
)

type timer struct {
	duration time.Duration
	start    *time.Time
	elapsed  time.Duration
}

var _ timr.Timer = (*timer)(nil)

func Timer(d time.Duration) *timer {
	return &timer{d, nil, 0}
}

func (t *timer) Resume(now time.Time) {
	if t.start == nil {
		t.start = &now
	}
}

func (t *timer) Pause(now time.Time) {
	if t.start != nil {
		t.elapsed += now.Sub(*t.start)
		t.start = nil
	}
}

func (t *timer) Reset() {
	t.start, t.elapsed = nil, 0
}

func (t *timer) Remaining(now time.Time) (remaining time.Duration, isRunning bool) {
	remaining = t.duration - t.elapsed

	if t.start == nil {
		return remaining, false
	}

	return remaining - now.Sub(*t.start), true
}
