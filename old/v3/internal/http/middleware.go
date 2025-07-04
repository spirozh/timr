package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

// isDevelopment checks if we're running in development mode
func isDevelopment() bool {
	return os.Getenv("GO_ENV") != "production"
}

// NoPanic is a middleware that recovers from panics and returns a 500 error response
// In development mode, it includes the stack trace in the response
func NoPanic(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic, URL, and stack trace
				log.Printf("panic on %s: %v\n%s", r.URL.String(), err, debug.Stack())

				// Ensure headers haven't been written before setting status
				if w.Header().Get("Content-Type") == "" {
					w.Header().Set("Content-Type", "application/json")
				}

				w.WriteHeader(http.StatusInternalServerError)

				// Prepare response based on environment
				response := map[string]any{
					"error": fmt.Sprintf("Internal Server Error: %s", err),
				}

				// Include stack trace in development mode
				if isDevelopment() {
					response["stack_trace"] = string(debug.Stack())
					response["error_details"] = err
				}

				// Best effort to write the error response
				if err := json.NewEncoder(w).Encode(response); err != nil {
					log.Printf("failed to write error response: %v", err)
				}
			}
		}()

		h.ServeHTTP(w, r)
	})
}
