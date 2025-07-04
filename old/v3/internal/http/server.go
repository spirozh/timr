package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type server struct {
	ctx             context.Context
	cancel          func()
	addr            string
	h               http.Handler
	shutdownTimeout time.Duration
}

func Server(ctx context.Context) server {
	ctx, done := context.WithCancel(ctx)

	return server{
		ctx:    ctx,
		cancel: done,
	}
}

func (s *server) Serve() error {
	srv := &http.Server{
		Addr:        s.addr,
		Handler:     s.h,
		BaseContext: func(_ net.Listener) context.Context { return s.ctx },
	}
	closingErrChan := make(chan error)
	go listenAndServe(srv, s.cancel, closingErrChan)

	return errors.Join(waitForShutdown(s.ctx, s.cancel, srv, s.shutdownTimeout), <-closingErrChan)
}

func listenAndServe(srv *http.Server, cancel func(), errChan chan<- error) {
	err := srv.ListenAndServe()
	if err == http.ErrServerClosed {
		err = nil
	}
	if err != nil {
		err = fmt.Errorf("listenAndServe error:\n%w", err)
	}
	cancel()
	errChan <- err
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
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("waitForShutdown error:\n%w", err)
	}

	return nil
}
