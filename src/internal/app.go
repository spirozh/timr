package internal

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type App struct {
	baseCtx  context.Context
	cancelFn func()
	err      chan error
}

func NewApp() App {
	ctx, cancel := context.WithCancel(context.Background())
	return App{
		baseCtx:  ctx,
		cancelFn: cancel,
		err:      make(chan error),
	}
}

func (a App) Run() {
	ctx := a.baseCtx

	mux := http.NewServeMux()
	AddRoutes(mux)

	s := &http.Server{
		BaseContext:       func(net.Listener) context.Context { return ctx },
		Addr:              ":8080",
		Handler:           LogRequest(mux),
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       4 * time.Second,
	}

	// Start the web server
	go func() {
		fmt.Println("Server listening on port 8080")
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for context to be closed and then shut down
	go func() {
		<-ctx.Done()
		a.err <- s.Shutdown(context.WithoutCancel(ctx))
	}()
}

func (a App) Close() error {
	a.cancelFn()
	return <-a.err
}
