package timr_test

import (
	"fmt"
	"spirozh/timr"
	"testing"
	"time"
)

func Test_TimerElapsed(t *testing.T) {
	var zero, now time.Time
	testTimer := timr.New("foo")

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

func Test_TimerEncoding(t *testing.T) {
	testTimer := timr.New("timer")
	var now time.Time

	t.Run("stopped timer", func(t *testing.T) {
		encodedTimer := testTimer.ToString()
		decodedTimer := timr.FromString(encodedTimer)
		if testTimer.Elapsed(now) != decodedTimer.Elapsed(now) {
			t.Fatal("timers don't really decode right")
		}
	})

	t.Run("running timer", func(t *testing.T) {
		testTimer.Start(now)

		encodedTimer := testTimer.ToString()
		fmt.Println(encodedTimer)
		decodedTimer := timr.FromString(encodedTimer)
		fmt.Println(decodedTimer.ToString())

		now = now.Add(time.Minute)
		if testTimer.Elapsed(now) != decodedTimer.Elapsed(now) {
			t.Fatal("timers don't really decode right")
		}
	})
}
