package http

import (
	"fmt"
	"net/http"

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

func APIRoutes(ts timr.TimerService) http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("/api/create/", nothingYet) // /api/create/name/duration
	m.HandleFunc("/api/list/", nothingYet)   // /api/list/
	m.HandleFunc("/api/pause/", nothingYet)  // /api/pause/name
	m.HandleFunc("/api/resume/", nothingYet) // /api/resume/name
	m.HandleFunc("/api/delete/", nothingYet) // /api/delete/name

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
