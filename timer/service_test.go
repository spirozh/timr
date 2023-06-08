package timer

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
	test.Equal(t, timr.ErrNoSuchTimer, s.Pause("x"))
	test.Equal(t, timr.ErrNoSuchTimer, s.Resume("x"))
	test.Equal(t, timr.ErrNoSuchTimer, s.Reset("x"))

	_, _, err = s.Remaining("x")
	test.Equal(t, timr.ErrNoSuchTimer, err)

	test.Equal(t, timr.ErrNoSuchTimer, s.Remove("x"))

	// when timer is exists there are no errors
	_ = s.Create("x", time.Minute)

	test.Equal(t, nil, s.Pause("x"))
	test.Equal(t, nil, s.Resume("x"))
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
	s.Resume("x")

	// recreating existing timer which is running and not yet expired fails
	test.Equal(t, timr.ErrTimerExists, s.Create("x", time.Minute))

	// even if it is almost expired
	now = now.Add(time.Minute)
	test.Equal(t, timr.ErrTimerExists, s.Create("x", time.Minute))

	// even if it is expired
	now = now.Add(time.Nanosecond)
	err = s.Create("x", time.Minute)
	test.Equal(t, timr.ErrTimerExists, err)

	// even if it is paused
	s.Create("x", time.Minute)
	err = s.Create("x", time.Minute)
	test.Equal(t, timr.ErrTimerExists, err)

	//t.Fail()
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
	ss := ts.(timr.Subscribable)

	var (
		calledEventType timr.ServiceEventType
		calledName      string
		calledCount     int
	)

	reset := func() {
		calledEventType, calledName, calledCount = 0, "", 0
	}

	ensure := func(call func(), eventType timr.ServiceEventType, name string, count int) {
		call()
		t.Helper()
		test.Equal(t, eventType, calledEventType)
		test.Equal(t, name, calledName)
		test.Equal(t, count, calledCount)
		reset()
	}

	cb := func(eventType timr.ServiceEventType, name string, _ timr.Timer) {
		calledEventType, calledName, calledCount = eventType, name, calledCount+1
	}

	sub1 := ss.Subscribe(cb)

	// when a call fails, there should be no notification
	ensure(func() { test.Equal(t, timr.ErrNoSuchTimer, ts.Pause("a")) }, 0, "", 0)
	ensure(func() { test.Equal(t, timr.ErrNoSuchTimer, ts.Resume("a")) }, 0, "", 0)
	ensure(func() { test.Equal(t, timr.ErrNoSuchTimer, ts.Reset("a")) }, 0, "", 0)
	ensure(func() { test.Equal(t, timr.ErrNoSuchTimer, ts.Remove("a")) }, 0, "", 0)

	// all events are emitted
	ensure(func() { ts.Create("a", 0) }, timr.EventTimerCreated, "a", 1)
	ensure(func() { ts.Resume("a") }, timr.EventTimerResumed, "a", 1)
	ensure(func() { ts.Pause("a") }, timr.EventTimerPaused, "a", 1)
	ensure(func() { ts.Reset("a") }, timr.EventTimerReset, "a", 1)
	ensure(func() { ts.Remove("a") }, timr.EventTimerRemoved, "a", 1)

	// double subscription (on the same callback) means double notification
	sub2 := ss.Subscribe(cb)
	ensure(func() { ts.Create("a", 0) }, timr.EventTimerCreated, "a", 2)
	ensure(func() { ts.Resume("a") }, timr.EventTimerResumed, "a", 2)
	ensure(func() { ts.Pause("a") }, timr.EventTimerPaused, "a", 2)
	ensure(func() { ts.Reset("a") }, timr.EventTimerReset, "a", 2)
	ensure(func() { ts.Remove("a") }, timr.EventTimerRemoved, "a", 2)

	ss.Unsubscribe(sub1)
	ss.Unsubscribe(sub2)
}

func TestTimerServiceUnsubscription(t *testing.T) {
	ts := TimerService(Now)
	ss := ts.(timr.Subscribable)

	// subscriptions are properly removed
	var s string
	suba := ss.Subscribe(func(e timr.ServiceEventType, n string, _ timr.Timer) { s += "a" })
	subb := ss.Subscribe(func(e timr.ServiceEventType, n string, _ timr.Timer) { s += "b" })
	subc := ss.Subscribe(func(e timr.ServiceEventType, n string, _ timr.Timer) { s += "c" })
	ts.Create("a", 0)
	test.Equal(t, "abc", s)

	ss.Unsubscribe(subb)
	ts.Reset("a")
	test.Equal(t, "abcac", s)

	ss.Unsubscribe(suba)
	ts.Reset("a")
	test.Equal(t, "abcacc", s)

	ss.Unsubscribe(subc)
	ts.Remove("a")
	test.Equal(t, "abcacc", s)
}
