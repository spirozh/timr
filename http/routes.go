package http

import (
	"net/http"

	"github.com/spirozh/timr/html"
)

func TemplateRoutes() http.Handler {
	m := http.NewServeMux()
	return m
}

func APIRoutes() http.Handler {
	m := http.NewServeMux()

	return m
}

func SSERoutes() http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("/api/sse/long/", Long)

	return m
}

func routes() http.Handler {
	m := http.NewServeMux()

	// templates
	m.Handle("/", TemplateRoutes())

	// static
	m.Handle("/static/", http.FileServer(http.FS(html.Static)))

	// api
	m.Handle("/api/", APIRoutes())

	// sse
	m.Handle("/api/sse/", SSERoutes())

	// selma
	m.HandleFunc("/selma", Selma)
	m.HandleFunc("/selma/", Selma)

	return m
}
