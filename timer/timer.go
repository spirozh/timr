package timer

import (
	"time"

	"github.com/spirozh/timr"
)

type timer struct {
	duration time.Duration
	start    *time.Time
	segments []time.Duration
}

var _ timr.Timer = (*timer)(nil)

func Timer(d time.Duration) *timer {
	return &timer{d, nil, nil}
}

func (t *timer) Resume(now time.Time) {
	if t.start == nil {
		t.start = &now
	}
}

func (t *timer) Pause(now time.Time) {
	if t.start != nil {
		t.segments = append(t.segments, now.Sub(*t.start))
		t.start = nil
	}
}

func (t *timer) Reset() {
	t.start, t.segments = nil, nil
}

func (t *timer) Remaining(now time.Time) (duration time.Duration, isRunning bool) {
	duration += t.duration
	for _, seg := range t.segments {
		duration -= seg
	}

	if t.start == nil {
		return duration, false
	}

	return duration - now.Sub(*t.start), true
}
