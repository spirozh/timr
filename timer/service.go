package timer

import (
	"time"

	"github.com/spirozh/timr"
	//"golang.org/x/exp/slices"
)

type timerService struct {
	clock  func() time.Time
	timers map[string]timr.Timer

	subscribers []*timr.EventSubscription
}

func TimerService(clock func() time.Time) timr.TimerService {
	ts := &timerService{
		clock:       clock,
		timers:      map[string]timr.Timer{},
		subscribers: []*timr.EventSubscription{},
	}

	return ts
}

func (ts *timerService) Subscribe(callback timr.EventCallback) *timr.EventSubscription {
	sub := &timr.EventSubscription{Callback: callback}
	ts.subscribers = append(ts.subscribers, sub)
	return sub
}

func (ts *timerService) subscriberIndex(sub *timr.EventSubscription) int {
	for i, subscriber := range ts.subscribers {
		if subscriber == sub {
			return i
		}
	}
	return -1
}

func (ts *timerService) Unsubscribe(sub *timr.EventSubscription) {
	// find the index of the subscription
	i := ts.subscriberIndex(sub)
	if i == -1 {
		return
	}

	// swap with the last one and slice the end off
	ts.subscribers[i], ts.subscribers[len(ts.subscribers)-1] = ts.subscribers[len(ts.subscribers)-1], ts.subscribers[i]
	ts.subscribers = ts.subscribers[:len(ts.subscribers)-1]
}

func (ts *timerService) notify(t timr.TimrEventType, name string, timer timr.Timer) {
	for _, sub := range ts.subscribers {
		sub.Callback(t, name, timer)
	}
}

func (ts *timerService) Create(name string, duration time.Duration) error {
	_, exists := ts.timers[name]
	if exists {
		return timr.ErrTimerExists
	}

	t := &timer{
		clock: ts.clock,
		notify: func(e timr.TimrEventType, t timr.Timer) {
			ts.notify(e, name, t)
		},
		duration: duration,
	}
	ts.timers[name] = t

	ts.notify(timr.Created, name, t)
	return nil
}

func (ts *timerService) Get(name string) (timr.Timer, error) {
	timer, ok := ts.timers[name]
	if ok {
		return timer, nil
	}

	return nil, timr.ErrNoSuchTimer
}

func (ts *timerService) Remove(name string) error {
	if _, ok := ts.timers[name]; !ok {
		return timr.ErrNoSuchTimer
	}

	delete(ts.timers, name)

	ts.notify(timr.Removed, name, nil)
	return nil
}

func (ts *timerService) List() []string {
	names := []string{}

	for name := range ts.timers {
		names = append(names, name)
	}

	return names
}
