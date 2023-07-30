package timr

import (
	"encoding/json"
	"time"
)

type timer struct {
	name    string
	start   *time.Time
	elapsed time.Duration
}

type timerExport struct {
	Name    string `json:"name"`
	Start   string `json:"start"`
	Elapsed int    `json:"elapsed"`
}

const timeFormat = "2006-01-02T15:04:05.0000"

func New(name string) timer {
	return timer{name: name}
}

func fromExport(te timerExport) timer {
	t := timer{
		name:    te.Name,
		elapsed: time.Duration(te.Elapsed) * time.Millisecond,
	}

	if te.Start != "" {
		start, err := time.Parse(timeFormat, te.Start)
		if err == nil {
			t.start = &start
		}
	}
	return t
}
func FromString(encoded string) timer {
	var te timerExport
	err := json.Unmarshal([]byte(encoded), &te)
	if err != nil {
		panic("just cannot")
	}
	return fromExport(te)
}

func (t timer) export() timerExport {
	te := timerExport{
		Name:    t.name,
		Elapsed: int(t.elapsed.Milliseconds()),
	}
	if t.start != nil {
		te.Start = t.start.Format(timeFormat)
	}
	return te
}

func (t timer) ToString() string {
	ex := t.export()
	b, _ := json.Marshal(ex)
	return string(b)
}

func (t timer) Elapsed(now time.Time) time.Duration {
	e := t.elapsed
	if t.start != nil {
		e += now.Sub(*t.start)
	}
	return e
}

func (t *timer) Start(now time.Time) {
	if t.start == nil {
		t.start = &now
	}
}

func (t *timer) Stop(now time.Time) {
	if t.start != nil {
		t.elapsed += now.Sub(*t.start)
		t.start = nil
	}
}
