package internal

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/spirozh/timr/internal/view"
)

func (a App) AddRoutes(mux *http.ServeMux) {
	t, err := view.NewTemplates()
	if err != nil {
		panic(err)
	}
	i := view.NewIndexView(t)

	mux.HandleFunc("GET /", i.Index)
	// Define the handler for the '/selma' route
	mux.HandleFunc("GET /selma", selma)
	mux.HandleFunc("GET /kill", a.kill)
}

func selma(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintln(w, "Hello, Selma!")
	if err != nil {
		slog.Error("Error writing response", "error", err)
	}
}

func (a App) kill(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintln(w, "Goodbye")
	if err != nil {
		slog.Error("Error writing response", "error", err)
	}
	a.Close()
}
