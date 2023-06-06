package memory

import (
	"testing"
	"time"

	"github.com/spirozh/timr"
	"github.com/spirozh/timr/test"
)

var now time.Time = time.Now()

func Now() time.Time {
	return now
}

func TestTimerRemaining(t *testing.T) {
	type testCase struct {
		label     string
		remaining time.Duration
		expired   bool
		duration  time.Duration
		start     *time.Time
		segments  []time.Duration
	}

	var t0 time.Time

	tc := func(label string, remaining time.Duration, expired bool, duration time.Duration, start *time.Time, elapsedSegments []time.Duration) {
		t.Helper()
		testTimer := timer{duration, start, elapsedSegments}

		format := "Test Case: '%s'"

		test.Equal(t, remaining, testTimer.remaining(t0), format, label)
		test.Equal(t, expired, testTimer.expired(t0), format, label)
	}

	noElapsedSegments := []time.Duration{}
	tc("dur 0m, not started, 0 segments", 0, false, 0, nil, noElapsedSegments)

	t0 = time.Now()
	tc("dur 0m, started t0, 0 segments", 0, false, 0, &t0, noElapsedSegments)

	tMinus1 := t0.Add(-time.Minute)
	tc("dur 0m, 1 min ago, 0 segments", -time.Minute, true, 0, &tMinus1, noElapsedSegments)

	tc("dur 4m, started now, 0 segments", 4*time.Minute, false, 4*time.Minute, &t0, noElapsedSegments)
	tc("dur 4m, not running, 0 segments", 4*time.Minute, false, 4*time.Minute, nil, noElapsedSegments)
	tc("dur 4m, 1 min ago, 0 segments", 3*time.Minute, false, 4*time.Minute, &tMinus1, noElapsedSegments)

	oneElapsedSegment := []time.Duration{time.Minute}
	tc("dur 0m, not running, 1 segment", -time.Minute, false, 0, nil, oneElapsedSegment)
	tc("dur 0m, 1 min ago, 1 segment", -2*time.Minute, false, 0, &tMinus1, oneElapsedSegment)

	tc("dur 4m, not running, 1 segment", 3*time.Minute, false, 4*time.Minute, nil, oneElapsedSegment)
	tc("dur 4m, 1 min ago, 1 segment", 2*time.Minute, false, 4*time.Minute, &tMinus1, oneElapsedSegment)

	tc("dur 4m, 1 min ago, 3 segments", 0, false, 4*time.Minute, &tMinus1, []time.Duration{2 * time.Minute, 30 * time.Second, 30 * time.Second})
}

func TestTimerService(t *testing.T) {
	s := TimerService(Now)
	test.NotEqual(t, s, nil, "TimerService(time.Time) should not return nil")
}

func TestTimerServiceNoSuchTimerError(t *testing.T) {
	var err error
	s := TimerService(Now)

	// no such timer errors
	test.Equal(t, timr.ErrNoSuchTimer, s.Toggle("x"))
	test.Equal(t, timr.ErrNoSuchTimer, s.Reset("x"))

	_, _, err = s.Remaining("x")
	test.Equal(t, timr.ErrNoSuchTimer, err)

	test.Equal(t, timr.ErrNoSuchTimer, s.Remove("x"))

	// when timer is exists there are no errors
	_ = s.Create("x", time.Minute)

	test.Equal(t, nil, s.Toggle("x"))
	test.Equal(t, nil, s.Reset("x"))

	_, _, err = s.Remaining("x")
	test.Equal(t, nil, err)

	test.Equal(t, nil, s.Remove("x"))
}

func TestTimerServiceCreateRunningTimerError(t *testing.T) {
	var err error

	s := TimerService(Now)

	// create timer and start it
	test.Equal(t, nil, s.Create("x", time.Minute))
	s.Toggle("x")

	// recreating existing timer which is running and not yet expired fails
	test.Equal(t, timr.ErrTimerRunning, s.Create("x", time.Minute))

	// remaining time of zero is still not expired
	now = now.Add(time.Minute)
	test.Equal(t, timr.ErrTimerRunning, s.Create("x", time.Minute))

	// recreating existing timer which is running but expired succeeds
	now = now.Add(time.Nanosecond)
	err = s.Create("x", time.Minute)
	test.Equal(t, nil, err)

	// recreating existing timer which is not running and not yet expired succeeds
	s.Create("x", time.Minute)
	err = s.Create("x", time.Minute)
	test.Equal(t, nil, err)
}

func TestTimerServiceListAndRemove(t *testing.T) {

	s := TimerService(Now)

	test.ElementsMatch(t, []string{}, s.List())

	names := []string{"a", "b", "c"}
	for _, name := range names {
		s.Create(name, 0)
	}

	test.ElementsMatch(t, names, s.List())

	for _, name := range names {
		s.Remove(name)
	}

	test.ElementsMatch(t, []string{}, s.List())
}

func TestTimerServiceRemaining(t *testing.T) {

}
