package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func Serve(ctx context.Context, cancel func(), shutdownTimeout time.Duration, addr string, h http.Handler) {
	srv := &http.Server{Addr: addr, Handler: h}

	var wg sync.WaitGroup
	wg.Add(1)

	go listenAndServe(srv, &wg)

	if err := waitForShutdown(ctx, cancel, srv, shutdownTimeout); err != nil {
		log.Println(fmt.Errorf("error shutting down server:\n%w", err))
	}

	wg.Wait()
}

func listenAndServe(srv *http.Server, wg *sync.WaitGroup) {
	defer wg.Done()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Println(fmt.Errorf("error closing server:\n%w", err))
	}
}

func waitForShutdown(ctx context.Context, ctxCancel func(), srv *http.Server, timeout time.Duration) error {
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal or the context is canceled
	select {
	case <-c:
		ctxCancel() // cancel the context, because an interrupt occurred
	case <-ctx.Done():
	}

	// Create a deadline to wait for server shutdown (use a new context, since the old context has already been cancelled).
	ctx, ctxCancel = context.WithTimeout(context.Background(), timeout)
	defer ctxCancel()

	// Waits for all connections to be closed, or until the timeout deadline
	return srv.Shutdown(ctx)
}
