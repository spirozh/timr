package timer

import (
	"encoding/json"
	"time"
)

type timerAux struct {
	Name     string `json:"name"`
	Duration int    `json:"duration"`
	Started  string `json:"started"`
	Elapsed  int    `json:"elapsed"`
}

const timeFormat = "2006-01-02T15:04:05.0000"

func (t timer) MarshalJSON() ([]byte, error) {
	started := ""
	if t.started != nil {
		started = t.started.Format(timeFormat)
	}
	aux := &timerAux{
		Name:     t.Name(),
		Duration: int(t.duration.Milliseconds()),
		Started:  started,
		Elapsed:  int(t.elapsed.Milliseconds()),
	}
	return json.Marshal(aux)
}

func (t *timer) UnmarshalJSON(data []byte) error {
	aux := &timerAux{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var started *time.Time
	if aux.Started != "" {
		auxTime, err := time.Parse(timeFormat, aux.Started)
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
