package memory

import (
	"time"
)

type timer struct {
	duration time.Duration
	start    *time.Time
	segments []time.Duration
}

func Timer(d time.Duration) *timer {
	return &timer{d, nil, nil}
}

func (t *timer) toggle(now time.Time) {
	if t.start != nil {
		t.segments = append(t.segments, now.Sub(*t.start))
		t.start = nil
	} else {
		t.start = &now
	}
}

func (t *timer) remaining(now time.Time) time.Duration {
	dur := t.duration
	for _, seg := range t.segments {
		dur -= seg
	}

	if t.start != nil {
		dur -= now.Sub(*t.start)
	}

	return dur
}

func (t *timer) expired(now time.Time) bool {
	return t.remaining(now) < 0
}
