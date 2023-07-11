package timr

import (
	"log"
	"os"
)

type TimerState struct {
	Duration  *int  `json:"duration"`
	Running   *bool `json:"running"`
	Remaining *int  `json:"remaining"`
}

type TimerMessage struct {
	Name  string     `json:"name"`
	State TimerState `json:"state"`
}

type Timer interface {
	State() TimerState
	Set(state TimerState)
}

type TimerService interface {
	Create(name string, state TimerState) (id int, err error)
	Update(id int, name string, state TimerState) error
	ForAll(func(id int, name string, state TimerState))
	Get(id int) (string, Timer, error)
	Remove(id int) error

	Subscribe(callback EventCallback) *EventSubscription
	Unsubscribe(*EventSubscription)
}

type EventSubscription struct {
	Callback EventCallback
}

type EventCallback func(eventType TimrEventType, id int, name string, timer Timer)

type TimrEventType int

const (
	_ TimrEventType = iota
	TimrEventCreated
	TimrEventSet
	TimrEventRemoved
)

func (t TimrEventType) String() string {
	switch t {
	case TimrEventCreated:
		return "TimrEventCreated"
	case TimrEventSet:
		return "TimrEventSet"
	case TimrEventRemoved:
		return "TimrEventRemoved"
	default:
		return "<unknown TimrEventType>"
	}
}

// errors
type timrError string

func (e timrError) Error() string {
	return string(e)
}

const (
	ErrNoSuchTimer timrError = "No Such Timer"
)

// logging

var (
	infoLog = log.New(os.Stderr, "INFO: ", log.Ltime)
	INFO    = func(info ...any) {
		infoLog.Println(info...)
	}
)
