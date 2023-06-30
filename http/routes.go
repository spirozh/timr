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
	prefix := "/"

	FileServer(m, prefix, html.FS)
	Selma(m, prefix)
	APIRoutes(m, prefix, ts) // everything under "/api"

	return m
}

func Selma(m *http.ServeMux, prefix string) {
	prefix += "selma/"
	timr.INFO("registering Selma at:\t\t", prefix)
	m.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello selma!!\n")
	})
}

func FileServer(m *http.ServeMux, prefix string, fsys fs.FS) {
	timr.INFO("registering FileServer at:\t\t", prefix)

	index, err := http.FS(fsys).Open("index.html")
	if err != nil {
		panic("can't open index.html")
	}
	info, err := index.Stat()
	if err != nil {
		panic("cant stat index.html")
	}

	fileServer := http.FileServer(http.FS(fsys))

	m.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeContent(w, r, "index.html", info.ModTime(), index)
		} else {
			fileServer.ServeHTTP(w, r)
		}
	})
}
