package internal

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (a App) AddRoutes(mux *http.ServeMux) {
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
