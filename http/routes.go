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
	fileServer := http.FileServer(http.FS(html.FS))
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fileServer.ServeHTTP(w, r) })

	// register api routes
	APIRoutes(m, "/api/timer/", ts) // register api routes for timer service

	return m
}
