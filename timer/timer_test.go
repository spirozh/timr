package timer

import (
	"testing"
	"time"

	"github.com/spirozh/timr"
	"github.com/spirozh/timr/test"
)

func TestTimerPauseResume(t *testing.T) {
	// TODO: think about this...
	mmss := func(s string) time.Time {
		t, _ := time.Parse("04:05", s)
		return t
	}

	dur := func(s string) time.Duration {
		return mmss(s).Sub(mmss("00:00"))
	}

	str := func(d time.Duration) string {
		f := "04:05"
		if d < 0 {
			d, f = -d, "-"+f
		}
		return mmss("00:00").Add(d).Format(f)
	}

	now := time.Now()
	clock := func() time.Time { return now }
	notify := func(timr.TimrEventType, timr.Timer) {}
	timer := &timer{clock, notify, dur("01:00"), nil, 0}

	advance := func(s string) {
		now = now.Add(dur(s))
	}

	shouldRemain := func(s string) {
		t.Helper()
		remaining, _ := timer.Remaining()
		if str(remaining) != s {
			t.Errorf("time remaining should be: %#v, was: %#v", s, str(remaining))
		}
	}

	timer.Resume()
	advance("00:05")
	timer.Pause()
	shouldRemain("00:55")
	advance("00:05")
	shouldRemain("00:55")
	timer.Resume()
	advance("00:55")
	shouldRemain("00:00")
	timer.Pause()
	advance("01:00")
	shouldRemain("00:00")
	timer.Resume()
	advance("01:00")
	shouldRemain("-01:00")

}

func TestTimerRemaining(t *testing.T) {
	t0 := time.Now()
	tMinus1 := t0.Add(-time.Minute)

	tc := func(remaining time.Duration, duration time.Duration, start *time.Time, elapsedTime time.Duration) {
		t.Helper()
		now := func() time.Time { return t0 }
		notify := func(e timr.TimrEventType, thisTimer timr.Timer) {}
		testTimer := timer{now, notify, duration, start, elapsedTime}

		actualRemaining, isRunning := testTimer.Remaining()
		test.Equal(t, remaining, actualRemaining)
		test.Equal(t, start != nil, isRunning)
	}

	tc(0, 0, nil, 0)                                        // dur 0m, not started, 0 segments
	tc(0, 0, &t0, 0)                                        // dur 0m, started t0, 0 segments
	tc(-time.Minute, 0, &tMinus1, 0)                        // dur 0m, 1 min ago, 0 segments
	tc(4*time.Minute, 4*time.Minute, &t0, 0)                // dur 4m, started t0, 0 segments
	tc(4*time.Minute, 4*time.Minute, nil, 0)                // dur 4m, not running, 0 segments
	tc(3*time.Minute, 4*time.Minute, &tMinus1, 0)           // dur 4m, 1 min ago, 0 segments
	tc(-2*time.Minute, 0, &tMinus1, time.Minute)            // dur 0m, 1 min ago, 1 segment
	tc(3*time.Minute, 4*time.Minute, nil, time.Minute)      // dur 4m, not running, 1 segment
	tc(2*time.Minute, 4*time.Minute, &tMinus1, time.Minute) // dur 4m, 1 min ago, 1 segment
	tc(0, 4*time.Minute, &tMinus1, 3*time.Minute)           // dur 4m, 1 min ago, 3 segments
}
