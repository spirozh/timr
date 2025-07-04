// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"spirozh/timr/internal"
	"syscall"
	"time"
)

func main() {
	baseCtx := context.Background()

	mux := http.NewServeMux()
	internal.AddRoutes(mux)

	// Create a new server instance
	s := &http.Server{
		BaseContext:       func(net.Listener) context.Context { return baseCtx },
		Addr:              ":8080",
		Handler:           internal.LogRequest(mux),
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       4 * time.Second,
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
