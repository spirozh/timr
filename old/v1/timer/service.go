package timer

import (
	"time"

	"github.com/spirozh/timr"
	//"golang.org/x/exp/slices"
)

type timerService struct {
	clock  func() time.Time
	lastID int
	timers []timerEntry

	subscribers []*timr.EventSubscription
}
type timerEntry struct {
	id    int
	name  string
	timer timr.Timer
}

func TimerService(clock func() time.Time) timr.TimerService {
	ts := &timerService{
		clock:       clock,
		timers:      []timerEntry{},
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

func (ts *timerService) notify(t timr.TimrEventType, id int, name string, timer timr.Timer) {
	for _, sub := range ts.subscribers {
		sub.Callback(t, id, name, timer)
	}
}

func (ts *timerService) Create(name string, state timr.TimerState) (int, error) {
	ts.lastID++

	var id = ts.lastID

	t := &timer{
		clock: ts.clock,
		notify: func(e timr.TimrEventType, t timr.Timer) {
			ts.notify(e, id, name, t)
		},
		duration: time.Millisecond * time.Duration(*state.Duration),
	}

	if state.Remaining != nil {
		t.elapsed = t.duration - time.Millisecond*time.Duration(*state.Remaining)
	}

	if state.Running != nil && *state.Running {
		now := ts.clock()
		t.start = &now
	}

	ts.timers = append(ts.timers, timerEntry{ts.lastID, name, t})

	ts.notify(timr.TimrEventCreated, id, name, t)
	return ts.lastID, nil
}

func (ts *timerService) Update(id int, name string, state timr.TimerState) error {
	// TODO: update timer ...
	return nil
}

func (ts *timerService) ForAll(f func(id int, name string, state timr.TimerState)) {
	for _, te := range ts.timers {
		f(te.id, te.name, te.timer.State())
	}
}

func (ts *timerService) Get(id int) (string, timr.Timer, error) {
	for _, te := range ts.timers {
		if te.id == id {
			return te.name, te.timer, nil
		}
	}

	return "", nil, timr.ErrNoSuchTimer
}

func (ts *timerService) Remove(id int) error {
	for i, te := range ts.timers {
		if te.id == id {
			for ; i < len(ts.timers)-1; i++ {
				ts.timers[i] = ts.timers[i+1]
			}
			ts.timers = ts.timers[:len(ts.timers)-1]

			ts.notify(timr.TimrEventRemoved, id, te.name, nil)
			return nil
		}
	}

	return timr.ErrNoSuchTimer
}
