package memory

import (
	"fmt"
	"strings"
	"time"

	"github.com/spirozh/timr"
)

type timer struct {
	duration time.Duration
	start    *time.Time
	segments []time.Duration
}

var zeroTimer *timer = &timer{}

func (t *timer) remaining(now time.Time) time.Duration {
	dur := t.duration
	for _, seg := range t.segments {
		dur -= seg
	}

	if t.start != nil {
		dur -= now.Sub(*t.start)
	}

	return dur
}

func (t *timer) expired(now time.Time) bool {
	return t.remaining(now) < 0
}

func (t timer) String() string {
	b := strings.Builder{}
	fmt.Fprint(&b, "<timer")
	if t.start == nil {
		fmt.Fprintf(&b, " %s", t.remaining(time.Now()))
	} else {
		fmt.Fprintf(&b, " %s *", t.remaining(time.Now()))
	}
	fmt.Fprint(&b, ">")
	return b.String()
}

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

	return zeroTimer, timr.ErrNoSuchTimer
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
	var names []string

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
	if timer.start != nil {
		timer.segments = append(timer.segments, now.Sub(*timer.start))
		timer.start = nil
	} else {
		timer.start = &now
	}

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
