package timer

import (
	"testing"
	"time"
)

func TestTimerRemaining(t *testing.T) {
	t0 := time.Now().UTC()
	tMinus1 := t0.Add(-time.Minute)

	tc := func(testname string, remaining time.Duration, duration time.Duration, start *time.Time, elapsed time.Duration) {
		t.Helper()
		t.Run(testname, func(t *testing.T) {
			t.Helper()
			clock := func() time.Time { return t0 }
			testTimer := timer{clock, nil, duration, start, elapsed}

			timerState := testTimer.State()
			dur := time.Duration(*timerState.Remaining) * time.Millisecond
			if remaining != dur {
				t.Errorf("state.Remaining: expected %v, got %v", remaining, dur)
			}

			if start == nil {
				if *timerState.Running {
					t.Error("state.Running: expected false")
				}
			} else {
				if !*timerState.Running {
					t.Error("state.Running: expected true")
				}
			}
		})
	}

	tc("dur 0m, not running, 0 segments", 0, 0, nil, 0)
	tc("dur 0m, started t0, 0 segments", 0, 0, &t0, 0)
	tc("dur 0m, 1 min ago, 0 segments", -time.Minute, 0, &tMinus1, 0)
	tc("dur 4m, started t0, 0 segments", 4*time.Minute, 4*time.Minute, &t0, 0)
	tc("dur 4m, not running, 0 segments", 4*time.Minute, 4*time.Minute, nil, 0)
	tc("dur 4m, 1 min ago, 0 segments", 3*time.Minute, 4*time.Minute, &tMinus1, 0)
	tc("dur 0m, 1 min ago, 1 segment", -2*time.Minute, 0, &tMinus1, time.Minute)
	tc("dur 4m, not running, 1 segment", 3*time.Minute, 4*time.Minute, nil, time.Minute)
	tc("dur 4m, 1 min ago, 1 segment", 2*time.Minute, 4*time.Minute, &tMinus1, time.Minute)
	tc("dur 4m, 1 min ago, 3 segments", 0, 4*time.Minute, &tMinus1, 3*time.Minute)
}

func TestTimerSet(t *testing.T) {
}

func TestTimerNotifications(t *testing.T) {
}
