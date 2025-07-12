package internal

import (
	"log/slog"
	"net/http"
	"time"
)

type CustomResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (w *CustomResponseWriter) WriteHeader(status int) {
	// receive status code from this method
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

func LogRequest(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		lrw := &CustomResponseWriter{w, http.StatusOK}
		next.ServeHTTP(lrw, r)
		slog.Info("request completed", "status", lrw.Status, "elapsed", time.Since(now), "method", r.Method, "path", r.URL.Path)
	}
}
