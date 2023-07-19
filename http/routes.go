package http

import (
	"io"
	"net/http"

	"github.com/spirozh/timr"
	"github.com/spirozh/timr/html"
)

func routes(ts timr.TimerService, sseDone chan any) http.Handler {
	m := http.NewServeMux()

	// register selma
	m.HandleFunc("/selma/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			io.WriteString(w, "hello selma!!\n")
			return
		case http.MethodHead:
			return
		default:
			w.Header().Add("Allow", "GET, HEAD")
			w.WriteHeader(405)
		}
	})

	// register file server
	m.Handle("/", http.FileServer(http.FS(html.FS)))

	// register api routes
	apiMux := http.NewServeMux()
	m.Handle("/api/", apiMux)

	apiMux.Handle("/api/timer/", newTimerHandler("/api/timer/", ts))
	apiMux.HandleFunc("/api/timer/sse", SSE(ts, sseDone))

	return m
}
