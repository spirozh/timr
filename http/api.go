package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/spirozh/timr"
)

// Selma handles /selma
func Selma(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello selma!!\n")
}

func resetTimer(m *http.ServeMux, prefix string, ts timr.TimerService) {
	prefix += "reset/"
	fmt.Println("registering resetTimer  at:", prefix)

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
	fmt.Println("registering deleteTimer at:", prefix)

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
	fmt.Println("registering sse         at:", prefix)
	m.HandleFunc(prefix, SSE(ts))
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
				fmt.Fprintf(w, "data: {\"%s\": null}\n\n", name)
			}

			flusher.Flush()
		}
		sub := ts.Subscribe(tsEventHandler)
		defer ts.Unsubscribe(sub)

		<-r.Context().Done()

		fmt.Print("EXITING SSE")
	}
}
