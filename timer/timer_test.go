package timer

import (
	"testing"
	"time"

	"github.com/spirozh/timr/test"
)

func TestTimerToggle(t *testing.T) {
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
	timer := Timer(dur("01:00")) // off

	advance := func(s string) {
		now = now.Add(dur(s))
	}
	shouldRemain := func(s string) {
		t.Helper()
		remaining, _ := timer.Remaining(now)
		if str(remaining) != s {
			t.Errorf("time remaining should be: %#v, was: %#v", s, str(remaining))
		}
	}

	timer.Resume(now)
	advance("00:05")
	timer.Pause(now)
	shouldRemain("00:55")
	advance("00:05")
	shouldRemain("00:55")
	timer.Resume(now)
	advance("00:55")
	shouldRemain("00:00")
	timer.Pause(now)
	advance("01:00")
	shouldRemain("00:00")
	timer.Resume(now)
	advance("01:00")
	shouldRemain("-01:00")

}

func TestTimerRemaining(t *testing.T) {
	t0 := time.Now()
	tMinus1 := t0.Add(-time.Minute)

	tc := func(remaining time.Duration, duration time.Duration, start *time.Time, elapsedSegments []time.Duration) {
		t.Helper()
		testTimer := timer{duration, start, elapsedSegments}

		actualRemaining, isRunning := testTimer.Remaining(t0)
		test.Equal(t, remaining, actualRemaining)
		test.Equal(t, start != nil, isRunning)
	}

	noElapsedSegments := []time.Duration{}
	oneElapsedSegment := []time.Duration{time.Minute}
	threeSegments := []time.Duration{2 * time.Minute, 30 * time.Second, 30 * time.Second}

	tc(0, 0, nil, noElapsedSegments)                              // dur 0m, not started, 0 segments
	tc(0, 0, &t0, noElapsedSegments)                              // dur 0m, started t0, 0 segments
	tc(-time.Minute, 0, &tMinus1, noElapsedSegments)              // dur 0m, 1 min ago, 0 segments
	tc(4*time.Minute, 4*time.Minute, &t0, noElapsedSegments)      // dur 4m, started t0, 0 segments
	tc(4*time.Minute, 4*time.Minute, nil, noElapsedSegments)      // dur 4m, not running, 0 segments
	tc(3*time.Minute, 4*time.Minute, &tMinus1, noElapsedSegments) // dur 4m, 1 min ago, 0 segments
	tc(-time.Minute, 0, nil, oneElapsedSegment)                   // dur 0m, not running, 1 segment
	tc(-2*time.Minute, 0, &tMinus1, oneElapsedSegment)            // dur 0m, 1 min ago, 1 segment
	tc(3*time.Minute, 4*time.Minute, nil, oneElapsedSegment)      // dur 4m, not running, 1 segment
	tc(2*time.Minute, 4*time.Minute, &tMinus1, oneElapsedSegment) // dur 4m, 1 min ago, 1 segment
	tc(0, 4*time.Minute, &tMinus1, threeSegments)                 // dur 4m, 1 min ago, 3 segments
}
