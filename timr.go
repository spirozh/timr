package timr

import (
	"time"
)

type Timer interface {
	Remaining() (remaining time.Duration, isPaused bool)
	Pause()
	Resume()
	Reset()
}

type TimerService interface {
	Create(name string, duration time.Duration) error
	List() []string
	Get(name string) (Timer, error)
	Remove(name string) error
}

type TimrEventType int

const (
	_ TimrEventType = iota
	EventTimerCreated
	EventTimerPaused
	EventTimerResumed
	EventTimerReset
	EventTimerRemoved

	timrEventNames string = "UnknownCreatedPausedResumedResetRemoved"
)

var timrEventNameOffsets = [...]int{0, 7, 14, 20, 27, 32, 39}

func (t TimrEventType) String() string {
	return timrEventNames[timrEventNameOffsets[t]:timrEventNameOffsets[t+1]]
}

type EventSubscription struct {
	Callback EventCallback
}

type EventCallback func(eventType TimrEventType, name string, timer Timer)

type Subscribable interface {
	Subscribe(callback EventCallback) *EventSubscription
	Unsubscribe(*EventSubscription)
}

// errors
type timrError string

func (e timrError) Error() string {
	return string(e)
}

const (
	ErrTimerExists timrError = "Timer Exists"
	ErrNoSuchTimer timrError = "No Such Timer"
)
