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

// Long handles /sse/long/
func Long(w http.ResponseWriter, r *http.Request) {
	// Make sure that the writer supports flushing.
	//
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	for i := 0; ; i++ {

		fmt.Fprintf(w, "data: %d\n\n", i)
		flusher.Flush()

		time.Sleep(time.Second)
	}
}
