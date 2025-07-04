// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func logRequest(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func selma(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "Hello Selma")
	if err != nil {
		log.Println("Error writing response:", err)
	}
}

func main() {
	baseCtx := context.Background()

	mux := http.NewServeMux()

	// Define the handler for the '/selma' route
	mux.HandleFunc("/selma", selma)

	// Create a new server instance
	s := &http.Server{
		Addr:    ":8080",
		Handler: logRequest(mux), // Use the default handler
	}

	// Start the web server in a goroutine
	go func() {
		fmt.Println("Server listening on port 8080")
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Set up signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(baseCtx, 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
