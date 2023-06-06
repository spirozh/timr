package memory

import (
	"testing"
	"time"

	"github.com/spirozh/timr"
	"github.com/stretchr/testify/assert"
)

type timerRemainingTestCase struct {
	label string
	rem   time.Duration
	dur   time.Duration
	st    *time.Time
	segs  []time.Duration
}

var now time.Time = time.Now()

func Now() time.Time {
	return now
}

func TestTimerRemaining(t *testing.T) {

	tcs := []timerRemainingTestCase{}

	tc := func(label string, remaining, duration time.Duration, start *time.Time, elapsedSegments []time.Duration) {
		tcs = append(tcs, timerRemainingTestCase{label, remaining, duration, start, elapsedSegments})
	}

	noElapsedSegments := []time.Duration{}
	tc("dur 0m, not started, 0 segs", 0, 0, nil, noElapsedSegments)

	t0 := time.Now()
	tc("dur 0m, started t0, 0 segs ", 0, 0, &t0, noElapsedSegments)

	tMinus1 := t0.Add(-time.Minute)
	tc("dur 0m, 1 min ago, 0 segs  ", -time.Minute, 0, &tMinus1, noElapsedSegments)

	tc("dur 4m, started now, 0 segs", 4*time.Minute, 4*time.Minute, &t0, noElapsedSegments)
	tc("dur 4m, not running, 0 segs", 4*time.Minute, 4*time.Minute, nil, noElapsedSegments)
	tc("dur 4m, 1 min ago, 0 segs  ", 3*time.Minute, 4*time.Minute, &tMinus1, noElapsedSegments)

	oneElapsedSegment := []time.Duration{time.Minute}
	tc("dur 0m, not running, 1 seg", -time.Minute, 0, nil, oneElapsedSegment)
	tc("dur 0m, 1 min ago, 1 seg  ", -2*time.Minute, 0, &tMinus1, oneElapsedSegment)

	tc("dur 4m, not running, 1 seg", 3*time.Minute, 4*time.Minute, nil, oneElapsedSegment)
	tc("dur 4m, 1 min ago, 1 seg  ", 2*time.Minute, 4*time.Minute, &tMinus1, oneElapsedSegment)

	tc("dur 4m, 1 min ago, 3 seg", 0, 4*time.Minute, &tMinus1, []time.Duration{2 * time.Minute, 30 * time.Second, 30 * time.Second})

	for _, tc := range tcs {
		testTimer := timer{tc.dur, tc.st, tc.segs}

		remaining := testTimer.remaining(t0)
		assert.Equal(t, tc.rem, remaining, "remaining: %s: %s", tc.label, testTimer)

		expired := testTimer.expired(t0)
		assert.Equal(t, tc.rem < 0, expired, "expired: %s: %s", tc.label, testTimer)
	}
}

func TestTimerService(t *testing.T) {
	s := TimerService(Now)
	assert.NotNil(t, s)
}

func TestTimerServiceNoSuchTimerError(t *testing.T) {
	var err error
	s := TimerService(Now)

	// no such timer errors
	err = s.Toggle("x")
	assert.ErrorIs(t, err, timr.ErrNoSuchTimer)

	err = s.Reset("x")
	assert.ErrorIs(t, err, timr.ErrNoSuchTimer)

	_, _, err = s.Remaining("x")
	assert.ErrorIs(t, err, timr.ErrNoSuchTimer)

	err = s.Remove("x")
	assert.ErrorIs(t, err, timr.ErrNoSuchTimer)

	// create timer
	_ = s.Create("x", time.Minute)

	// when timer is exists there are no errors
	err = s.Toggle("x")
	assert.NoError(t, err)

	err = s.Reset("x")
	assert.NoError(t, err)

	_, _, err = s.Remaining("x")
	assert.NoError(t, err)

	err = s.Remove("x")
	assert.NoError(t, err)
}

func TestTimerServiceCreateRunningTimerError(t *testing.T) {
	var err error

	s := TimerService(Now)

	// create timer
	err = s.Create("x", time.Minute)
	assert.NoError(t, err)

	s.Toggle("x")

	// recreating existing timer which is running and not yet expired fails
	err = s.Create("x", time.Minute)
	assert.ErrorIs(t, err, timr.ErrTimerRunning)

	// remaining time of zero is still not expired
	now = now.Add(time.Minute)
	err = s.Create("x", time.Minute)
	assert.ErrorIs(t, err, timr.ErrTimerRunning)

	// recreating existing timer which is running but expired succeeds
	now = now.Add(time.Nanosecond)
	err = s.Create("x", time.Minute)
	assert.NoError(t, err)

	// recreating existing timer which is not running and not yet expired succeeds
	s.Create("x", time.Minute)
	err = s.Create("x", time.Minute)
	assert.NoError(t, err)
}

func TestTimerServiceListAndRemove(t *testing.T) {

	s := TimerService(Now)

	assert.ElementsMatch(t, []string{}, s.List())

	names := []string{"a", "b", "c"}
	for _, name := range names {
		s.Create(name, 0)
	}

	assert.ElementsMatch(t, names, s.List())

	for _, name := range names {
		s.Remove(name)
	}

	assert.ElementsMatch(t, []string{}, s.List())
}

func TestTimerServiceRemaining(t *testing.T) {

}
