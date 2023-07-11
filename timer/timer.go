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

func (t *timer) Set(state timr.TimerState) {
	if state.Running != nil {
		if *state.Running {
			if t.start == nil {
				now := t.clock()
				t.start = &now
			}
		} else {
			if t.start != nil {
				t.elapsed += t.clock().Sub(*t.start)
				t.start = nil
			}
		}
	}

	if state.Duration != nil && state.Remaining == nil {
		t.duration = ms(*state.Duration)
	}
	if state.Duration == nil && state.Remaining != nil {
		t.elapsed = t.duration - ms(*state.Remaining)

		if t.start != nil {
			now := t.clock()
			t.start = &now
		}
	}
	if state.Duration != nil && state.Remaining != nil {
		t.duration = ms(*state.Duration)
		t.elapsed = t.duration - ms(*state.Remaining)

		if t.start != nil {
			now := t.clock()
			t.start = &now
		}
	}

	if t.notify != nil {
		t.notify(timr.TimrEventSet, t)
	}
}

func ms(ms int) time.Duration {
	return time.Duration(ms)
}

func (t *timer) State() timr.TimerState {
	var duration, remaining int
	var running bool

	duration = int(t.duration.Milliseconds())
	remaining = duration - int(t.elapsed.Milliseconds())

	if t.start != nil {
		running = true
		remaining -= int(t.clock().Sub(*t.start).Milliseconds())
	}

	return timr.TimerState{
		Duration:  &duration,
		Remaining: &remaining,
		Running:   &running,
	}
}
