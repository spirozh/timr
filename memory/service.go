package memory

import (
	"time"

	"github.com/spirozh/timr"
	"github.com/spirozh/timr/timer"
	"golang.org/x/exp/slices"
)

type timerService struct {
	clock  func() time.Time
	timers map[string]timr.Timer

	subscribers []*timr.EventSubscription
}

var _ timr.Subscribable = (*timerService)(nil)

func TimerService(clock func() time.Time) timr.TimerService {
	return &timerService{
		clock:  clock,
		timers: map[string]timr.Timer{},

		subscribers: []*timr.EventSubscription{},
	}
}

func (s *timerService) Subscribe(callback timr.EventCallback) *timr.EventSubscription {
	sub := &timr.EventSubscription{Callback: callback}
	s.subscribers = append(s.subscribers, sub)
	return sub
}

func (s *timerService) Unsubscribe(sub *timr.EventSubscription) {
	// find the index of the subscription
	i := slices.Index(s.subscribers, sub)
	if i == -1 {
		return
	}

	// swap withthe last one and reslice
	s.subscribers[i], s.subscribers[len(s.subscribers)-1] = s.subscribers[len(s.subscribers)-1], s.subscribers[i]
	s.subscribers = s.subscribers[:len(s.subscribers)-1]
}

func (s *timerService) notify(t timr.ServiceEventType, name string) {
	for _, sub := range s.subscribers {
		sub.Callback(t, name)
	}
}

func (s *timerService) getTimer(name string) (timr.Timer, error) {
	timer, ok := s.timers[name]
	if ok {
		return timer, nil
	}

	return nil, timr.ErrNoSuchTimer
}

func (s *timerService) Create(name string, duration time.Duration) error {
	_, exists := s.timers[name]
	if exists {
		return timr.ErrTimerExists
	}

	s.timers[name] = timer.Timer(duration)

	s.notify(timr.EventTimerCreated, name)
	return nil
}

func (s *timerService) List() []string {
	names := []string{}

	for name := range s.timers {
		names = append(names, name)
	}

	return names
}

func (s *timerService) Pause(name string) error {
	timer, err := s.getTimer(name)
	if err != nil {
		return err
	}

	timer.Pause(s.clock())

	s.notify(timr.EventTimerPaused, name)
	return nil
}

func (s *timerService) Resume(name string) error {
	timer, err := s.getTimer(name)
	if err != nil {
		return err
	}

	timer.Resume(s.clock())

	s.notify(timr.EventTimerResumed, name)
	return nil
}

func (s *timerService) Reset(name string) error {
	timer, err := s.getTimer(name)
	if err != nil {
		return err
	}

	timer.Reset()

	s.notify(timr.EventTimerReset, name)
	return nil
}

func (s *timerService) Remaining(name string) (remaining time.Duration, isPaused bool, err error) {
	timer, err := s.getTimer(name)
	if err != nil {
		return 0, false, err
	}

	remaining, paused := timer.Remaining(s.clock())
	return remaining, paused, nil
}

func (s *timerService) Remove(name string) error {
	_, err := s.getTimer(name)
	if err != nil {
		return err
	}

	delete(s.timers, name)

	s.notify(timr.EventTimerRemoved, name)
	return nil
}
