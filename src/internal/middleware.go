package internal

import (
	"log/slog"
	"net/http"
)

func LogRequest(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	}
}
