package timr

import (
	"time"
)

type TimerService interface {
	Create(name string, duration time.Duration) error
	List() []string
	Toggle(name string) error
	Reset(name string) error
	Remaining(name string) (remaining time.Duration, isRunning bool, err error)
	Remove(name string) error
}

type ServiceEventType int

const (
	EventTimerCreated ServiceEventType = iota
	EventTimerToggled
	EventTimerReset
	EventTimerRemove
)

type EventSubscription struct {
}

type XXX interface {
	Subscribe(selector string, callback func(timerName string, eventType ServiceEventType)) *EventSubscription
	Unsubscribe(*EventSubscription) error
}

// errors
type timrError string

func (e timrError) Error() string {
	return string(e)
}

const ErrTimerRunning timrError = "Timer Running"
const ErrNoSuchTimer timrError = "No Such Timer"
