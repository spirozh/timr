package http

import (
	"io"
	"io/fs"
	"net/http"

	"github.com/spirozh/timr"
	"github.com/spirozh/timr/html"
)

func routes(ts timr.TimerService) http.Handler {
	m := http.NewServeMux()

	FileServer(m, "/", html.FS)
	Selma(m, "/selma/")
	APIRoutes(m, "/api/", ts) // everything under "/api"

	return m
}

func FileServer(m *http.ServeMux, path string, fsys fs.FS) {
	timr.INFO("registering FileServer at:\t\t", path)

	fileServer := http.FileServer(http.FS(fsys))

	m.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	})
}

func Selma(m *http.ServeMux, path string) {
	timr.INFO("registering Selma at:\t\t", path)

	m.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello selma!!\n")
	})
}
