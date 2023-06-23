package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/spirozh/timr"
)

// Selma handles /selma
func Selma(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello selma!!\n")
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
