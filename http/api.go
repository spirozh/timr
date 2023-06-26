package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spirozh/timr"
)

func apiRoutes(ts timr.TimerService) http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("/api/create/", createTimer(ts)) // /api/create/name/duration
	m.HandleFunc("/api/pause/", pauseTimer(ts))   // /api/pause/name
	m.HandleFunc("/api/resume/", resumeTimer(ts)) // /api/resume/name
	m.HandleFunc("/api/delete/", deleteTimer(ts)) // /api/delete/name

	// sse route
	m.HandleFunc("/api/sse", SSE(ts))

	return m
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

// SSE handles /api/sse/
func SSE(ts timr.TimerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("SSE!!")

		flusher, ok := sseSetup(w, r)
		if !ok {
			return
		}

		// first connection: send json for all timers
		m := map[string]timr.TimerState{}
		for _, name := range ts.List() {
			fmt.Printf("output data for %v\n", name)

			t, _ := ts.Get(name)
			m[name] = t.State()
			j, err := json.Marshal(m)
			if err == nil {
				fmt.Fprintf(w, "data: %s\n\n", string(j))
				flusher.Flush()
			}
			delete(m, name)
		}

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
				m[name] = timer.State()
				j, err := json.Marshal(m)
				if err == nil {
					fmt.Fprintf(w, "data: %s\n\n", string(j))
					flusher.Flush()
				}
				delete(m, name)

				return
			}

			fmt.Fprintf(w, "data: {\"%s\": null}\n\n", name)
			flusher.Flush()
		}
		sub := ts.Subscribe(tsEventHandler)
		defer ts.Unsubscribe(sub)

		<-r.Context().Done()

		fmt.Print("EXITING SSE")
	}
}

func createTimer(ts timr.TimerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		arg, ok := strings.CutPrefix(r.URL.Path, "/api/create/")
		args := strings.Split(arg, "/")
		if !ok || len(args) != 2 {
			fmt.Fprintf(w, "bad request: %s/n", r.URL.Path)
		}

		t, err := strconv.Atoi(args[1])
		if err == nil {
			err = ts.Create(args[0], time.Duration(t)*time.Millisecond)
		}

		fmt.Fprintf(w, "create: %v, err: %v\n", r.URL, err)
	}
}

func pauseTimer(ts timr.TimerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, ok := strings.CutPrefix(r.URL.Path, "/api/pause/")
		if !ok {
			fmt.Fprintf(w, "bad request: %s/n", r.URL.Path)
		}

		var err error
		if t, err := ts.Get(name); err == nil {
			t.Pause()
		}

		fmt.Fprintf(w, "pause: %v, err: %v\n", r.URL, err)
	}
}

func resumeTimer(ts timr.TimerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, ok := strings.CutPrefix(r.URL.Path, "/api/resume/")
		if !ok {
			fmt.Fprintf(w, "bad request: %s/n", r.URL.Path)
		}

		var err error
		if t, err := ts.Get(name); err == nil {
			t.Resume()
		}

		fmt.Fprintf(w, "resume: %v, err: %v\n", r.URL, err)
	}
}

func deleteTimer(ts timr.TimerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, ok := strings.CutPrefix(r.URL.Path, "/api/delete/")
		if !ok {
			fmt.Fprintf(w, "bad request: %s/n", r.URL.Path)
		}

		err := ts.Remove(name)

		fmt.Fprintf(w, "delete: %v, err: %v\n", r.URL, err)
	}
}
