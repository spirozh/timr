package timr

import (
	"time"
)

type Timer interface {
	Pause(time.Time)
	Resume(time.Time)
	Reset()
	Remaining(time.Time) (remaining time.Duration, isPaused bool)
}

type TimerService interface {
	Create(name string, duration time.Duration) error
	List() []string
	Pause(name string) error
	Resume(name string) error
	Reset(name string) error
	Remaining(name string) (remaining time.Duration, isPaused bool, err error)
	Remove(name string) error
}

type ServiceEventType int

const (
	_ ServiceEventType = iota
	EventTimerCreated
	EventTimerPaused
	EventTimerResumed
	EventTimerReset
	EventTimerRemoved

	serviceEventNames string = "UnknownCreatedPausedResumedResetRemoved"
)

var serviceEventNameOffsets = [...]int{0, 7, 14, 20, 27, 32, 39}

func (t ServiceEventType) String() string {
	return serviceEventNames[serviceEventNameOffsets[t]:serviceEventNameOffsets[t+1]]
}

type EventSubscription struct {
	Callback EventCallback
}

type EventCallback func(eventType ServiceEventType, name string, timer Timer)

type Subscribable interface {
	Subscribe(callback EventCallback) *EventSubscription
	Unsubscribe(*EventSubscription)
}

// errors
type timrError string

func (e timrError) Error() string {
	return string(e)
}

const ErrTimerExists timrError = "Timer Exists"
const ErrNoSuchTimer timrError = "No Such Timer"
