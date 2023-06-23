package timr

import (
	"time"
)

type TimerState struct {
	Running   bool  `json:"running"`
	Remaining int64 `json:"remaining"`
}

type Timer interface {
	State() TimerState
	Pause()
	Resume()
	Reset()
}

type TimerService interface {
	Create(name string, duration time.Duration) error
	List() []string
	Get(name string) (Timer, error)
	Remove(name string) error

	Subscribe(callback EventCallback) *EventSubscription
	Unsubscribe(*EventSubscription)
}

type EventSubscription struct {
	Callback EventCallback
}

type EventCallback func(eventType TimrEventType, name string, timer Timer)

type TimrEventType int

const (
	_ TimrEventType = iota
	Created
	Paused
	Resumed
	Reset
	Removed

	timrEventNames string = "UnknownCreatedPausedResumedResetRemoved"
)

var timrEventNameOffsets = [...]int{0, 7, 14, 20, 27, 32, 39}

func (t TimrEventType) String() string {
	return timrEventNames[timrEventNameOffsets[t]:timrEventNameOffsets[t+1]]
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
