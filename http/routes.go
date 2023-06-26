package http

import (
	"io"
	"net/http"

	"github.com/spirozh/timr"
	"github.com/spirozh/timr/html"
)

func routes(ts timr.TimerService) http.Handler {
	m := http.NewServeMux()

	// filesystem
	m.Handle("/", http.FileServer(http.FS(html.FS)))

	// api
	m.Handle("/api/", apiRoutes(ts))

	// selma
	m.HandleFunc("/selma", Selma)
	m.HandleFunc("/selma/", Selma)

	return m
}

// Selma handles /selma
func Selma(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello Selma!!\n")
}
