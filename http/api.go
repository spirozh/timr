package http

import (
	"fmt"
	"io"
	"net/http"
	"time"

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
	return
}

type TimerState struct {
	Running   bool  `json:"running"`
	Remaining int64 `json:"remaining"`
}

// SSE handles /api/sse/token
func SSE(ts timr.TimerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := sseSetup(w, r)
		if !ok {
			return
		}

		ticker := time.NewTicker(time.Second)

		// first connection: send json for all timers

		// after each state change, send json for just one timer
		//
		// either:
		//  {'name': {'running': bool, 'remaining': milliseconds}}
		// or (for delete)
		//  {'name': null}

		for i := 0; ; i++ {
			fmt.Printf("SSE: %d\n", i)
			fmt.Fprintf(w, "data: %d\n\n", i)
			flusher.Flush()

			select {
			case <-r.Context().Done():
				ticker.Stop()
				return
			case <-ticker.C:
				continue
			}
		}
	}
}
