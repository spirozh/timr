package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spirozh/timr"
)

func APIRoutes(parentM *http.ServeMux, prefix string, ts timr.TimerService) {
	prefix += "api/"
	timr.INFO("registering APIRoutes at:\t\t", prefix)

	m := http.NewServeMux()

	createTimer(m, prefix, ts) // /api/create/:name/:duration
	pauseTimer(m, prefix, ts)  // /api/pause/name
	resumeTimer(m, prefix, ts) // /api/resume/name
	resetTimer(m, prefix, ts)  // /api/reset/name
	deleteTimer(m, prefix, ts) // /api/delete/name
	sse(m, prefix, ts)         // /api/sse/

	parentM.Handle(prefix, m)
}

func createTimer(m *http.ServeMux, prefix string, ts timr.TimerService) {
	prefix += "create/"
	timr.INFO("registering + createTimer at:\t", prefix)

	m.HandleFunc(prefix,
		func(w http.ResponseWriter, r *http.Request) {
			arg, ok := strings.CutPrefix(r.URL.Path, prefix)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "misconfiguration: mux path should start with ", prefix)
				return
			}

			args := strings.Split(arg, "/")
			if !ok || len(args) != 2 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "b: bad request: %s (args: %v)\n", r.URL.Path, args)
				return
			}

			name := args[0]

			t, err := strconv.Atoi(args[1])
			if err != nil || t <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "c: bad request: %s (args: %v) %v\n", r.URL.Path, args, err)
				return
			}

			err = ts.Create(name, time.Duration(t)*time.Millisecond)
			if err == timr.ErrTimerExists {
				w.WriteHeader(http.StatusConflict)
				fmt.Fprintf(w, "d: bad request: %s, %v\n", r.URL.Path, err)
				return
			}

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "e: bad request: %s, name: %s, err: %v\n", r.URL.Path, name, err)
				return
			}

			fmt.Fprintf(w, "create: %v, name: %s, err: %v\n", r.URL, name, err)
		})
}

func pauseTimer(m *http.ServeMux, prefix string, ts timr.TimerService) {
	prefix += "pause/"
	timr.INFO("registering + pauseTimer at:\t", prefix)

	m.HandleFunc(prefix,
		func(w http.ResponseWriter, r *http.Request) {
			name, ok := strings.CutPrefix(r.URL.Path, prefix)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "misconfiguration: mux path should start with ", prefix)
				return
			}

			t, err := ts.Get(name)
			if err == nil {
				t.Pause()
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}

			fmt.Fprintf(w, "pause: %v, name: %s, err: %v\n", r.URL, name, err)
		})
}

func resumeTimer(m *http.ServeMux, prefix string, ts timr.TimerService) {
	prefix += "resume/"
	timr.INFO("registering + resumeTimer at:\t", prefix)

	m.HandleFunc(prefix,
		func(w http.ResponseWriter, r *http.Request) {
			name, ok := strings.CutPrefix(r.URL.Path, prefix)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "misconfiguration: mux path should start with ", prefix)
				return
			}

			t, err := ts.Get(name)
			if err == nil {
				t.Resume()
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}

			fmt.Fprintf(w, "resume: %v, name: %s, err: %v\n", r.URL, name, err)
		})
}

func resetTimer(m *http.ServeMux, prefix string, ts timr.TimerService) {
	prefix += "reset/"
	timr.INFO("registering + resetTimer at:\t", prefix)

	m.HandleFunc(prefix,
		func(w http.ResponseWriter, r *http.Request) {
			name, ok := strings.CutPrefix(r.URL.Path, prefix)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "misconfiguration: mux path should start with ", prefix)
				return
			}

			t, err := ts.Get(name)
			if err == nil {
				t.Reset()
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}

			fmt.Fprintf(w, "resume: %v, name: %s, err: %v\n", r.URL, name, err)
		})
}

func deleteTimer(m *http.ServeMux, prefix string, ts timr.TimerService) {
	prefix += "delete/"
	timr.INFO("registering + deleteTimer at:\t", prefix)

	m.HandleFunc(prefix,
		func(w http.ResponseWriter, r *http.Request) {
			name, ok := strings.CutPrefix(r.URL.Path, prefix)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "misconfiguration: mux path should start with ", prefix)
				return
			}

			err := ts.Remove(name)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}

			fmt.Fprintf(w, "delete: %v, name: %s, err: %v\n", r.URL, name, err)
		})
}

func sse(m *http.ServeMux, prefix string, ts timr.TimerService) {
	prefix += "sse/"
	timr.INFO("registering + sse at:\t\t", prefix)
	m.HandleFunc(prefix, SSE(ts))
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
		for _, name := range ts.List() {
			outputTimer(w, ts, name)
		}
		flusher.Flush()

		// after each state change, send json for just one timer
		//
		// either:
		//  {'name': {'running': bool, 'remaining': milliseconds}}
		// or (for delete)
		//  {'name': null}
		mu := sync.Mutex{}
		tsEventHandler := func(eventType timr.TimrEventType, name string, timer timr.Timer) {
			mu.Lock()
			defer mu.Unlock()

			if eventType != timr.Removed {
				outputTimer(w, ts, name)
			} else {
				outputNull(w, name)
			}

			flusher.Flush()
		}
		sub := ts.Subscribe(tsEventHandler)
		defer ts.Unsubscribe(sub)

		<-r.Context().Done()

		timr.INFO("ending SSE handler")
	}
}

func outputTimer(w io.Writer, ts timr.TimerService, name string) {
	t, _ := ts.Get(name)
	state := t.State()
	msg := timr.TimerMessage{Name: name, State: state}
	j, err := json.Marshal(msg)
	if err == nil {
		output(w, string(j))
	}
}

func outputNull(w io.Writer, name string) {
	output(w, fmt.Sprintf(`{"name":"%v"}`, name))
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
