package main

import (
	"context"
	"fmt"
	"log"

	// "spirozh/timr/internal/http"
	"time"
)

type app struct {
	rootContext     context.Context
	rootCancel      func()
	addr            string
	shutdownTimeout time.Duration
	// server          http.Server
}

func App() app {
	ctx, done := context.WithCancel(context.Background())

	return app{
		rootContext:     ctx,
		rootCancel:      done,
		addr:            ":8080",
		shutdownTimeout: time.Second,
	}
}

func (a *app) run() error {
	ctx, cancel := context.WithCancel(a.rootContext)
	defer cancel()

	// Initialize server
	if err := a.server.Init( /* config */ ); err != nil {
		return fmt.Errorf("server initialization failed: %w", err)
	}

	// Start server
	go func() {
		if err := a.server.Start(ctx); err != nil {
			log.Printf("server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
	defer shutdownCancel()

	return a.server.Shutdown(shutdownCtx)
}
