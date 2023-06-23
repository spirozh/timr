package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spirozh/timr"
	"github.com/spirozh/timr/html"
)

func TemplateRoutes() http.Handler {
	m := http.NewServeMux()
	return m
}

func nothingYet(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("Path: %s (not handled)\n", r.URL.Path)
	fmt.Print(msg)
	fmt.Fprint(w, msg)
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

func APIRoutes(ts timr.TimerService) http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("/api/create/", createTimer(ts)) // /api/create/name/duration
	m.HandleFunc("/api/pause/", pauseTimer(ts))   // /api/pause/name
	m.HandleFunc("/api/resume/", resumeTimer(ts)) // /api/resume/name
	m.HandleFunc("/api/delete/", deleteTimer(ts)) // /api/delete/name

	// sse route
	m.HandleFunc("/api/sse", SSE(ts))

	return m
}

func routes(ts timr.TimerService) http.Handler {
	m := http.NewServeMux()

	// filesystem
	m.Handle("/", http.FileServer(http.FS(html.FS)))

	// api
	m.Handle("/api/", APIRoutes(ts))

	// selma
	m.HandleFunc("/selma", Selma)
	m.HandleFunc("/selma/", Selma)

	return m
}
