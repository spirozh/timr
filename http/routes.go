package http

import (
	"io"
	"net/http"

	"github.com/spirozh/timr"
	"github.com/spirozh/timr/html"
)

func routes(ts timr.TimerService) http.Handler {
	m := http.NewServeMux()

	prefix := "/"

	// api
	m.Handle("/api/", apiRoutes(ts))

	// selma
	Selma(m, prefix)

	// filesystem
	FileServer(m, prefix, html.FS)

	return m
}

// Selma handles /selma
func Selma(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello Selma!!\n")
}
