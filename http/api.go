package http

import (
	"fmt"
	"io"
	"net/http"
	"time"
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

// SSE handles /api/sse/token
func SSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := sseSetup(w, r)
	if !ok {
		return
	}

	ticker := time.NewTicker(time.Second)

	// send json for all timers
	//  {'name': {'running': false, 'remaining': seconds}}
	//
	// send updates
	//  - new timer
	//  {'name': {'running': false, 'remaining': 2000}}
	//  - pause timer
	//  {'name': {'running': false}}
	//  - unpause timer
	//  {'name': {'running': true}}
	//  - delete timer
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
