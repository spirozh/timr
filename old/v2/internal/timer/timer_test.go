package timer_test

import (
	"fmt"
	"spirozh/timr/internal/timer"
	"testing"
	"time"
)

func Test_TimerElapsed(t *testing.T) {
	var zero, now time.Time
	testTimer := timer.New()

	testElapsed := func(expected time.Duration) {
		t.Helper()
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

func Test_TimerString(t *testing.T) {
	testCases := []struct {
		timr timer.Timer
		str  string
	}{
		{timer.New(timer.Name("foo")), `<timer "foo" duration:0s started:nil elapsed:0s>`},
		{timer.New(timer.Started(&time.Time{})), `<timer "" duration:0s started:"0001-01-01T00:00:00.0000" elapsed:0s>`},
	}
	for _, test := range testCases {
		t.Run(test.str, func(t *testing.T) {
			if actual := fmt.Sprint(test.timr); actual != test.str {
				t.Fatal(actual)
			}
		})
	}
}
