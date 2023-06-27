package http

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"

	"github.com/spirozh/timr"
	"github.com/spirozh/timr/html"
)

func TemplateRoutes() http.Handler {
	m := http.NewServeMux()
	return m
}

func routes(ts timr.TimerService) http.Handler {
	m := http.NewServeMux()

	prefix := "/"

	// api
	APIRoutes(m, prefix, ts) // everything under "/api"

	// selma
	Selma(m, prefix)

	// filesystem
	FileServer(m, prefix, html.FS)

	return m
}

func FileServer(m *http.ServeMux, prefix string, fsys fs.FS) {
	fmt.Println("registering FileServer  at:", prefix)
	m.Handle(prefix, http.FileServer(http.FS(fsys)))
}

// Selma handles /selma
func Selma(m *http.ServeMux, prefix string) {
	prefix += "selma/"
	fmt.Println("registering Selma       at:", prefix)
	m.HandleFunc(prefix, selma)
}

func selma(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello selma!!\n")
}
