package timer

import (
	"testing"
	"time"

	"github.com/spirozh/timr"
	"github.com/spirozh/timr/test"
)

var now time.Time = time.Now().UTC()

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
	_, err = s.Get("x")
	test.Equal(t, timr.ErrNoSuchTimer, err)
	test.Equal(t, timr.ErrNoSuchTimer, s.Remove("x"))

	// when timer is exists there are no errors
	_ = s.Create("x", time.Minute)

	_, err = s.Get("x")
	test.Equal(t, nil, err)
	test.Equal(t, nil, s.Remove("x"))
}

func TestTimerServiceCreateExistingTimerError(t *testing.T) {
	s := TimerService(Now)

	// create timer
	test.Equal(t, nil, s.Create("x", time.Minute))
	test.Equal(t, timr.ErrTimerExists, s.Create("x", time.Minute))
	test.Equal(t, nil, s.Remove("x"))
	test.Equal(t, nil, s.Create("x", time.Minute))
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

func TestTimerServiceSubscription(t *testing.T) {
	ts := TimerService(Now)

	var (
		calledEventType timr.TimrEventType
		calledName      string
		calledCount     int
	)

	ensure := func(call func(), eventType timr.TimrEventType, name string, count int) {
		call()
		t.Helper()
		test.Equal(t, eventType, calledEventType)
		test.Equal(t, name, calledName)
		test.Equal(t, count, calledCount)

		calledEventType, calledName, calledCount = 0, "", 0
	}

	cb := func(eventType timr.TimrEventType, name string, _ timr.Timer) {
		calledEventType, calledName, calledCount = eventType, name, calledCount+1
	}

	sub1 := ts.Subscribe(cb)

	// when a call fails, there should be no notification
	ensure(func() { test.Equal(t, timr.ErrNoSuchTimer, func() error { _, err := ts.Get("a"); return err }()) }, 0, "", 0)
	ensure(func() { test.Equal(t, timr.ErrNoSuchTimer, ts.Remove("a")) }, 0, "", 0)

	// all events are emitted
	ensure(func() { ts.Create("a", 0) }, timr.Created, "a", 1)
	ensure(func() { ts.Remove("a") }, timr.Removed, "a", 1)

	// double subscription (on the same callback) means double notification
	sub2 := ts.Subscribe(cb)
	ensure(func() { ts.Create("a", 0) }, timr.Created, "a", 2)
	ensure(func() { ts.Remove("a") }, timr.Removed, "a", 2)

	ts.Unsubscribe(sub1)
	ts.Unsubscribe(sub2)
}

func TestTimerServiceUnsubscription(t *testing.T) {
	ts := TimerService(Now)

	// subscriptions are properly removed
	var s string
	suba := ts.Subscribe(func(e timr.TimrEventType, n string, _ timr.Timer) { s += "a" })
	subb := ts.Subscribe(func(e timr.TimrEventType, n string, _ timr.Timer) { s += "b" })
	subc := ts.Subscribe(func(e timr.TimrEventType, n string, _ timr.Timer) { s += "c" })
	ts.Create("a", 0)
	test.Equal(t, "abc", s)

	ts.Unsubscribe(subb)
	ts.Remove("a")
	test.Equal(t, "abcac", s)

	ts.Unsubscribe(suba)
	ts.Create("a", 0)
	test.Equal(t, "abcacc", s)

	ts.Unsubscribe(subc)
	ts.Remove("a")
	test.Equal(t, "abcacc", s)
}
