package timer

import (
	"testing"
	"time"

	"github.com/spirozh/timr"
	"github.com/spirozh/timr/test"
)

func noNotify(timr.TimrEventType, timr.Timer) {}

func TestTimerReset(t *testing.T) {
	now := time.Now().UTC()
	clock := func() time.Time { return now }
	timer := &timer{clock, noNotify, dur("01:00"), nil, 0}

	timer.Resume()

}

func TestTimerPauseResumeReset(t *testing.T) {
	now := time.Now().UTC()
	clock := func() time.Time { return now }
	timer := &timer{clock, noNotify, dur("01:00"), nil, 0}

	timer.Resume()
	advance(&now, "00:05")
	timer.Pause()
	shouldRemain(t, timer, "00:55")
	advance(&now, "00:05")
	shouldRemain(t, timer, "00:55")
	timer.Resume()
	advance(&now, "00:55")
	shouldRemain(t, timer, "00:00")
	timer.Pause()
	advance(&now, "01:00")
	shouldRemain(t, timer, "00:00")
	timer.Resume()
	advance(&now, "01:00")
	shouldRemain(t, timer, "-01:00")
	timer.Reset()
	shouldRemain(t, timer, "01:00")
	advance(&now, "00:30")
	shouldRemain(t, timer, "00:30")
}

func TestTimerRemaining(t *testing.T) {
	t0 := time.Now().UTC()
	tMinus1 := t0.Add(-time.Minute)

	tc := func(remaining time.Duration, duration time.Duration, start *time.Time, elapsed time.Duration) {
		t.Helper()
		clock := func() time.Time { return t0 }
		testTimer := timer{clock, noNotify, duration, start, elapsed}

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

func TestTimerNotifications(t *testing.T) {
	var (
		calledEventType timr.TimrEventType
		calledTimer     *timer
	)

	ensure := func(call func(), e timr.TimrEventType, pt *timer) {
		call()
		t.Helper()
		test.Equal(t, e, calledEventType)
		test.Equal(t, pt, calledTimer)
		calledEventType, calledTimer = 0, nil
	}

	testTimer := &timer{time.Now, nil, time.Minute, nil, 0}
	ensure(testTimer.Pause, 0, nil)
	ensure(testTimer.Resume, 0, nil)
	ensure(testTimer.Reset, 0, nil)

	(*testTimer).notify = func(e timr.TimrEventType, t timr.Timer) {
		calledEventType, calledTimer = e, t.(*timer)
	}
	ensure(testTimer.Pause, timr.EventTimerPaused, testTimer)
	ensure(testTimer.Resume, timr.EventTimerResumed, testTimer)
	ensure(testTimer.Reset, timr.EventTimerReset, testTimer)
}

func mmss(s string) time.Time {
	t, _ := time.Parse("04:05", s)
	return t
}

func dur(s string) time.Duration {
	return mmss(s).Sub(mmss("00:00"))
}

func str(d time.Duration) string {
	f := "04:05"
	if d < 0 {
		d, f = -d, "-"+f
	}
	return mmss("00:00").Add(d).Format(f)
}

func advance(now *time.Time, s string) {
	*now = now.Add(dur(s))
}

func shouldRemain(t *testing.T, timer timr.Timer, s string) {
	t.Helper()
	remaining, _ := timer.Remaining()
	if str(remaining) != s {
		t.Errorf("time remaining should be: %#v, was: %#v", s, str(remaining))
	}
}
