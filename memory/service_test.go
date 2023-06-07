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

func TestTimerService(t *testing.T) {
	s := TimerService(Now)
	test.NotEqual(t, nil, s)
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
