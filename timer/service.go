package timer

import (
	"time"

	"github.com/spirozh/timr"
	"golang.org/x/exp/slices"
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

func (ts *timerService) Unsubscribe(sub *timr.EventSubscription) {
	// find the index of the subscription
	i := slices.Index(ts.subscribers, sub)
	if i == -1 {
		return
	}

	// swap withthe last one and reslice
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

	notify := func(e timr.TimrEventType, t timr.Timer) {
		ts.notify(e, name, t)
	}

	t := &timer{ts.clock, notify, duration, nil, 0}
	ts.timers[name] = t

	ts.notify(timr.EventTimerCreated, name, t)
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
	t, err := ts.Get(name)
	if err != nil {
		return err
	}

	delete(ts.timers, name)

	ts.notify(timr.EventTimerRemoved, name, t)
	return nil
}

func (ts *timerService) List() []string {
	names := []string{}

	for name := range ts.timers {
		names = append(names, name)
	}

	return names
}
