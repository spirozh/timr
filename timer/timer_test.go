package timer

import (
	"testing"
	"time"

	"github.com/spirozh/timr"
	"github.com/spirozh/timr/test"
)

func noNotify(timr.TimrEventType, timr.Timer) {}

func TestTimerRemaining(t *testing.T) {
	t0 := time.Now().UTC()
	tMinus1 := t0.Add(-time.Minute)

	tc := func(remaining time.Duration, duration time.Duration, start *time.Time, elapsed time.Duration) {
		t.Helper()
		clock := func() time.Time { return t0 }
		testTimer := timer{clock, noNotify, duration, start, elapsed}

		timerState := testTimer.State()
		test.Equal(t, remaining.Milliseconds(), timerState.Remaining)
		test.Equal(t, start != nil, timerState.Running)
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
	ensure(func() { testTimer.Set(timr.TimerState{}) }, 0, nil)

	(*testTimer).notify = func(e timr.TimrEventType, t timr.Timer) {
		calledEventType, calledTimer = e, t.(*timer)
	}
	ensure(func() { testTimer.Set(timr.TimerState{}) }, timr.TimrEventSet, testTimer)
}
