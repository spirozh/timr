package timer

import (
	"encoding/json"
	"time"
)

type timerAux struct {
	Name     string  `json:"name"`
	Duration int     `json:"duration"`
	Started  *string `json:"started"`
	Elapsed  int     `json:"elapsed"`
}

const timeFormat = "2006-01-02T15:04:05.0000"

func (t Timer) MarshalJSON() ([]byte, error) {
	var started *string
	if t.started != nil {
		timeStr := t.started.Format(timeFormat)
		started = &timeStr
	}
	aux := &timerAux{
		Name:     t.Name(),
		Duration: int(t.duration.Milliseconds()),
		Started:  started,
		Elapsed:  int(t.elapsed.Milliseconds()),
	}
	return json.Marshal(aux)
}

func (t *Timer) UnmarshalJSON(data []byte) error {
	aux := &timerAux{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var started *time.Time
	if aux.Started != nil {
		auxTime, err := time.Parse(timeFormat, *aux.Started)
		if err != nil {
			return err
		}
		started = &auxTime
	}

	t.Config(
		Name(aux.Name),
		Duration(time.Duration(aux.Duration)*time.Millisecond),
		Started(started),
		Elapsed(time.Duration(aux.Elapsed)*time.Millisecond),
	)

	return nil
}
