package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func Serve(ctx context.Context, h http.Handler) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	var mu sync.Mutex
	go listenAndServe(srv, &mu)
	waitForShutdown(ctx, srv)
	mu.Lock()
}

func listenAndServe(srv *http.Server, mu sync.Locker) {
	mu.Lock()
	defer mu.Unlock()

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func waitForShutdown(ctx context.Context, srv *http.Server) {
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal or context is canceled
	select {
	case <-c:
	case <-ctx.Done():
	}

	// Create a deadline to wait for server shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline.
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Println("server Shutdown error: ", err)
	}
}
