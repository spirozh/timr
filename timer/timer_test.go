package timer_test

import (
	"spirozh/timr/timer"
	"testing"
	"time"
)

func Test_TimerElapsed(t *testing.T) {
	var zero, now time.Time
	testTimer := timer.New()

	testElapsed := func(expected time.Duration) {
		elapsed := testTimer.Elapsed(now)
		if elapsed != expected {
			t.Fatal("expected: ", expected, ", was: ", elapsed)
		}
	}

	// 0s. zero elapsed time
	testElapsed(0)

	now = zero.Add(time.Second)
	testElapsed(0)       // 1s. zero elapsed time
	testTimer.Start(now) // 1s. start timer

	now = zero.Add(2 * time.Second)
	testElapsed(time.Second) // 2s. one second elapsed

	testTimer.Start(now) // 2s. multiple starts don't change anything
	testElapsed(time.Second)

	testTimer.Stop(now) // 2s. stop timer
	testElapsed(time.Second)

	now = zero.Add(3 * time.Second)
	testElapsed(time.Second) // 3s. stopped timer doesn't advance

	now = zero.Add(4 * time.Second)
	testTimer.Stop(now) // 4s. multiple stops don't change anything
	testElapsed(time.Second)
}
