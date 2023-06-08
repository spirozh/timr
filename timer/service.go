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

var _ timr.Subscribable = (*timerService)(nil)

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

func (ts *timerService) notify(t timr.ServiceEventType, name string, timer timr.Timer) {
	for _, sub := range ts.subscribers {
		sub.Callback(t, name, timer)
	}
}

func (ts *timerService) getTimer(name string) (timr.Timer, error) {
	timer, ok := ts.timers[name]
	if ok {
		return timer, nil
	}

	return nil, timr.ErrNoSuchTimer
}

func (ts *timerService) Create(name string, duration time.Duration) error {
	_, exists := ts.timers[name]
	if exists {
		return timr.ErrTimerExists
	}

	t := Timer(duration)
	ts.timers[name] = t

	ts.notify(timr.EventTimerCreated, name, t)
	return nil
}

func (ts *timerService) List() []string {
	names := []string{}

	for name := range ts.timers {
		names = append(names, name)
	}

	return names
}

func (ts *timerService) Pause(name string) error {
	timer, err := ts.getTimer(name)
	if err != nil {
		return err
	}

	timer.Pause(ts.clock())

	ts.notify(timr.EventTimerPaused, name, timer)
	return nil
}

func (ts *timerService) Resume(name string) error {
	timer, err := ts.getTimer(name)
	if err != nil {
		return err
	}

	timer.Resume(ts.clock())

	ts.notify(timr.EventTimerResumed, name, timer)
	return nil
}

func (ts *timerService) Reset(name string) error {
	timer, err := ts.getTimer(name)
	if err != nil {
		return err
	}

	timer.Reset()

	ts.notify(timr.EventTimerReset, name, timer)
	return nil
}

func (ts *timerService) Remaining(name string) (remaining time.Duration, isPaused bool, err error) {
	timer, err := ts.getTimer(name)
	if err != nil {
		return 0, false, err
	}

	remaining, paused := timer.Remaining(ts.clock())
	return remaining, paused, nil
}

func (ts *timerService) Remove(name string) error {
	_, err := ts.getTimer(name)
	if err != nil {
		return err
	}

	delete(ts.timers, name)

	ts.notify(timr.EventTimerRemoved, name, nil)
	return nil
}
