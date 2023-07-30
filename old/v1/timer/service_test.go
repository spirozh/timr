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
	_, _, err = s.Get(-1)
	test.Equal(t, timr.ErrNoSuchTimer, err)
	test.Equal(t, timr.ErrNoSuchTimer, s.Remove(1))

	// when timer is exists there are no errors
	id, _ := s.Create("x", durstate(1000))

	_, _, err = s.Get(id)
	test.Equal(t, nil, err)
	test.Equal(t, nil, s.Remove(id))
}

func TestTimerServiceForAll(t *testing.T) {
	s := TimerService(Now)

	count := 0
	s.ForAll(func(_ int, _ string, _ timr.TimerState) { count++ })
	test.Equal(t, 0, count)

	names := []string{"a", "b", "c"}
	for _, name := range names {
		s.Create(name, durstate(0))
	}

	s.ForAll(func(_ int, _ string, _ timr.TimerState) { count++ })
	test.Equal(t, 3, count)

	s.ForAll(func(id int, _ string, _ timr.TimerState) {
		err := s.Remove(id)
		test.Equal(t, nil, err)
	})

	count = 0
	s.ForAll(func(_ int, _ string, _ timr.TimerState) { count++ })
	test.Equal(t, 0, count)
}

func TestTimerServiceSubscribe(t *testing.T) {
	ts := TimerService(Now)

	var (
		calledEventType timr.TimrEventType
		calledId        int
		calledName      string
		calledCount     int
	)

	ensure := func(call func(), eventType timr.TimrEventType, id int, name string, count int) {
		call()
		t.Helper()
		test.Equal(t, eventType, calledEventType)
		test.Equal(t, id, calledId)
		test.Equal(t, name, calledName)
		test.Equal(t, count, calledCount)

		calledEventType, calledName, calledCount = 0, "", 0
	}

	cb := func(eventType timr.TimrEventType, id int, name string, _ timr.Timer) {
		calledEventType, calledId, calledName, calledCount = eventType, id, name, calledCount+1
	}

	sub1 := ts.Subscribe(cb)

	// when a call fails, there should be no notification
	ensure(func() { test.Equal(t, timr.ErrNoSuchTimer, func() error { _, _, err := ts.Get(1); return err }()) }, 0, 0, "", 0)
	ensure(func() { test.Equal(t, timr.ErrNoSuchTimer, ts.Remove(1)) }, 0, 0, "", 0)

	// all events are emitted
	var id int
	ensure(func() { id, _ = ts.Create("a", durstate(0)) }, timr.TimrEventCreated, id, "a", 1)
	// TODO: check that TimrEventSet is emitted
	ensure(func() { ts.Remove(id) }, timr.TimrEventRemoved, id, "a", 1)

	// double subscription (on the same callback) means double notification
	sub2 := ts.Subscribe(cb)
	ensure(func() { id, _ = ts.Create("a", durstate(0)) }, timr.TimrEventCreated, id, "a", 2)
	ensure(func() { ts.Remove(id) }, timr.TimrEventRemoved, id, "a", 2)

	ts.Unsubscribe(sub1)
	ts.Unsubscribe(sub2)
}

func durstate(ms int) timr.TimerState {
	return timr.TimerState{
		Duration: &ms}
}

func TestTimerServiceUnsubscribe(t *testing.T) {
	ts := TimerService(Now)

	// subscriptions are properly removed
	var s string
	subA := ts.Subscribe(func(e timr.TimrEventType, i int, n string, _ timr.Timer) { s += "a " })
	subB := ts.Subscribe(func(e timr.TimrEventType, i int, n string, _ timr.Timer) { s += "b " })
	subC := ts.Subscribe(func(e timr.TimrEventType, i int, n string, _ timr.Timer) { s += "c " })
	test.NotEqual(t, subA, subB)
	test.NotEqual(t, subA, subC)
	test.NotEqual(t, subB, subC)

	running := false
	one := 1
	state := timr.TimerState{
		Duration: &one,
		Running:  &running,
	}

	a, _ := ts.Create("a", state)
	test.Equal(t, "a b c ", s)

	ts.Unsubscribe(subB)
	ts.Remove(a)
	test.Equal(t, "a b c a c ", s)

	ts.Unsubscribe(subA)
	a, _ = ts.Create("a", state)
	test.Equal(t, "a b c a c c ", s)

	ts.Unsubscribe(subB)
	ts.Remove(a)
	test.Equal(t, "a b c a c c c ", s)

	ts.Unsubscribe(subC)
	a, _ = ts.Create("a", state)
	test.Equal(t, "a b c a c c c ", s)

	ts.Unsubscribe(subB)
	ts.Remove(a)
	test.Equal(t, "a b c a c c c ", s)

}
