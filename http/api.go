package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/spirozh/timr"
)

type TimerHandler struct {
	path string
	timr.TimerService
}

type timerMessage struct {
	Id    *int             `json:"id"`
	Name  *string          `json:"name"`
	State *timr.TimerState `json:"state"`
}

func APIRoutes(parentM *http.ServeMux, path string, ts timr.TimerService) {
	timr.INFO("registering APIRoutes at:\t\t", path)

	m := http.NewServeMux()

	m.HandleFunc(path+"sse", SSE(ts))
	m.Handle(path, TimerHandler{path, ts})

	parentM.Handle(path, m)
}

func (th TimerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// dispatch on METHOD
	switch r.Method {
	case http.MethodPost:
		// createTimer POST /api/timer/
		//
		// receives {"name":"zzz","state":{"duration":1,"remaining":1,"running":true}}
		// emits {"id":1,"name":"zzz","state":{"duration":1,"remaining":1,"running":true}}
		createReq, err := timerBody(w, r)
		if err != nil {
			return
		}
		if createReq.Id != nil {
			fmt.Println("TODO")
			// problem
		}
		if createReq.Name == nil {
			// problem
			fmt.Println("TODO")
		}

		id, _ := th.Create(*createReq.Name, *createReq.State)
		_, timer, _ := th.Get(id)
		s := timer.State()

		res, err := json.Marshal(timerMessage{&id, createReq.Name, &s})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// write something here about the error
			return
		}
		w.Write(res)
		return

	case http.MethodPatch:
		// updateTimer PATCH /api/timer/:id
		//
		// receives {"name":"mmm","state":{"remaining":5}}
		// emits {"id":1,"name":"mmm","state":{"duration":1,"remaining":5,"running":false}}
		id, _ := id(r)

		updateReq, err := timerBody(w, r)
		if err != nil {
			return
		}
		if updateReq.Id != nil && *updateReq.Id != id {
			// problem
			fmt.Print("TODO")
		}

		name, timer, _ := th.TimerService.Get(id)

		if updateReq.Name == nil {
			updateReq.Name = &name
		}

		if updateReq.State == nil {
			state := timer.State()
			updateReq.State = &state
		}

		// update the timer...
		th.TimerService.Update(id, *updateReq.Name, *updateReq.State)
		s := timer.State()

		res, err := json.Marshal(timerMessage{&id, updateReq.Name, &s})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// write something here about the error
			return
		}
		w.Write(res)
		return

	case http.MethodGet:
		if id, exists := id(r); !exists {
			// getTimers(m, prefix+"timer/")       // GET /api/timer/
			// emits [{"id":1,"name":"zzz","duration":1,"remaining":1},...]
		} else {
			// getTimer(m, prefix+"timer/")        // GET /api/timer/:id
			// get one timer
			fmt.Print(id)
		}

	case http.MethodDelete:
		// deleteTimer(m, prefix, ts)          // DELETE /api/timer/:id
		id, _ := id(r)
		fmt.Print(id)
	default:
		// 405
	}
}

func id(r *http.Request) (int, bool) {
	return 1, false
}

func timerBody(w http.ResponseWriter, r *http.Request) (timerMessage, error) {
	return timerMessage{}, nil
}

// SSE handles /api/sse/
func SSE(ts timr.TimerService) http.HandlerFunc {
	count := 0

	return func(w http.ResponseWriter, r *http.Request) {
		count++
		defer func() { count-- }()
		timr.INFO("starting SSE handler (", count, ")")

		flusher, ok := sseSetup(w, r)
		if !ok {
			return
		}

		// first connection: send json for all timers
		//m := map[string]timr.TimerState{}
		ts.ForAll(func(id int, name string, state timr.TimerState) {
			outputTimer(w, ts, id)
		})
		flusher.Flush()

		// after each state change, send json for just one timer
		//
		// either:
		//  {'name': {'running': bool, 'remaining': milliseconds}}
		// or (for delete)
		//  {'name': null}
		mu := sync.Mutex{}
		tsEventHandler := func(eventType timr.TimrEventType, id int, name string, timer timr.Timer) {
			mu.Lock()
			defer mu.Unlock()

			outputTimer(w, ts, id)

			flusher.Flush()
		}
		sub := ts.Subscribe(tsEventHandler)
		defer ts.Unsubscribe(sub)

		<-r.Context().Done()

		timr.INFO("ending SSE handler")
	}
}

func outputTimer(w io.Writer, ts timr.TimerService, id int) {
	name, t, _ := ts.Get(id)
	state := t.State()
	msg := timr.TimerMessage{Name: name, State: state}
	j, err := json.Marshal(msg)
	if err == nil {
		output(w, string(j))
	}
}

func output(w io.Writer, s string) {
	io.WriteString(w, "data: ")
	io.WriteString(w, s)
	io.WriteString(w, "\n\n")
}

func sseSetup(w http.ResponseWriter, r *http.Request) (flusher http.Flusher, ok bool) {
	// Make sure that the writer supports flushing.
	flusher, ok = w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	flusher.Flush()
	return
}
