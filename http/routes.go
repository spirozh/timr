package http

import (
	"io"
	"net/http"

	"github.com/spirozh/timr"
	"github.com/spirozh/timr/html"
)

func routes(ts timr.TimerService) http.Handler {
	m := http.NewServeMux()

	// register selma
	m.HandleFunc("/selma/", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "hello selma!!\n") })

	// register file server
	m.Handle("/", http.FileServer(http.FS(html.FS)))

	// register api routes
	apiMux := http.NewServeMux()
	m.Handle("/api/", apiMux)

	apiMux.Handle("/api/timer/", newTimerHandler("/api/timer/", ts))
	apiMux.HandleFunc("/api/timer/sse", SSE(ts))

	return m
}
