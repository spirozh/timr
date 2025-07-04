package internal

import (
	"fmt"
	"log/slog"
	"net/http"
)

func AddRoutes(mux *http.ServeMux) {
	// Define the handler for the '/selma' route
	mux.HandleFunc("GET /selma", selma)
}

func selma(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintln(w, "Hello Selma")
	if err != nil {
		slog.Error("Error writing response", "error", err)
	}
}
