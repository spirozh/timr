package memory

import (
	"time"

	"github.com/spirozh/timr"
)

type timerService struct {
	clock  func() time.Time
	timers map[string]*timer
}

func TimerService(clock func() time.Time) timr.TimerService {
	return &timerService{
		clock:  clock,
		timers: make(map[string]*timer),
	}
}

func (s *timerService) getTimer(name string) (*timer, error) {
	timer, ok := s.timers[name]
	if ok {
		return timer, nil
	}

	return nil, timr.ErrNoSuchTimer
}

func (s *timerService) Create(name string, duration time.Duration) error {
	now := s.clock()

	t, exists := s.timers[name]
	if exists && t.start != nil && !t.expired(now) {
		return timr.ErrTimerRunning
	}

	s.timers[name] = &timer{
		start:    nil,
		duration: duration,
	}
	return nil
}

func (s *timerService) List() []string {
	names := []string{}

	for name := range s.timers {
		names = append(names, name)
	}

	return names
}

func (s *timerService) Toggle(name string) error {
	timer, err := s.getTimer(name)
	if err != nil {
		return err
	}

	now := s.clock()
	timer.toggle(now)

	return nil
}

func (s *timerService) Reset(name string) error {
	timer, err := s.getTimer(name)
	if err != nil {
		return err
	}

	timer.segments = []time.Duration{}

	return nil
}

func (s *timerService) Remaining(name string) (remaining time.Duration, isRunning bool, err error) {
	timer, err := s.getTimer(name)
	if err != nil {
		return 0, false, err
	}

	return timer.remaining(s.clock()), timer.start != nil, nil
}

func (s *timerService) Remove(name string) error {
	_, err := s.getTimer(name)
	if err != nil {
		return err
	}

	delete(s.timers, name)

	return nil
}
