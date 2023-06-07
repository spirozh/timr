package memory

import (
	"testing"
	"time"
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

	gap := func(timer *timer, cur, pause, rem string) {
		t.Helper()

		timer.toggle(mmss(cur))
		timer.toggle(mmss(cur).Add(dur(pause)))
		act := str(timer.remaining(mmss(cur)))
		if rem != act {
			t.Errorf("%s should remain, but was %s", rem, act)
		}
	}

	timer := Timer(dur("01:00")) // off

	gap(timer, "00:01", "00:10", "00:50") // on
	gap(timer, "00:11", "00:00", "00:50")
	gap(timer, "00:01", "00:00", "00:50")
}

func TestTimerRemaining(t *testing.T) {
	t0 := time.Now()
	tMinus1 := t0.Add(-time.Minute)

	tc := func(remaining time.Duration, expired bool, duration time.Duration, start *time.Time, elapsedSegments []time.Duration) {
		t.Helper()
		testTimer := timer{duration, start, elapsedSegments}

		actualRemaining := testTimer.remaining(t0)
		if remaining != actualRemaining {
			t.Errorf("expected remaining=%v, was: %v", remaining, actualRemaining)
		}

		actualExpired := testTimer.expired(t0)
		if expired != actualExpired {
			t.Errorf("expected expired=%v, was: %v, ", expired, actualExpired)
		}
	}

	noElapsedSegments := []time.Duration{}
	oneElapsedSegment := []time.Duration{time.Minute}
	threeSegments := []time.Duration{2 * time.Minute, 30 * time.Second, 30 * time.Second}

	tc(0, false, 0, nil, noElapsedSegments)                              // dur 0m, not started, 0 segments
	tc(0, false, 0, &t0, noElapsedSegments)                              // dur 0m, started t0, 0 segments
	tc(-time.Minute, true, 0, &tMinus1, noElapsedSegments)               // dur 0m, 1 min ago, 0 segments
	tc(4*time.Minute, false, 4*time.Minute, &t0, noElapsedSegments)      // dur 4m, started t0, 0 segments
	tc(4*time.Minute, false, 4*time.Minute, nil, noElapsedSegments)      // dur 4m, not running, 0 segments
	tc(3*time.Minute, false, 4*time.Minute, &tMinus1, noElapsedSegments) // dur 4m, 1 min ago, 0 segments
	tc(-time.Minute, true, 0, nil, oneElapsedSegment)                    // dur 0m, not running, 1 segment
	tc(-2*time.Minute, true, 0, &tMinus1, oneElapsedSegment)             // dur 0m, 1 min ago, 1 segment
	tc(3*time.Minute, false, 4*time.Minute, nil, oneElapsedSegment)      // dur 4m, not running, 1 segment
	tc(2*time.Minute, false, 4*time.Minute, &tMinus1, oneElapsedSegment) // dur 4m, 1 min ago, 1 segment
	tc(0, false, 4*time.Minute, &tMinus1, threeSegments)                 // dur 4m, 1 min ago, 3 segments
}
