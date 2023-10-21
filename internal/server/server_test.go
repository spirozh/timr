package server_test

import (
	"context"
	"errors"
	"net/http"
	"spirozh/timr/internal/server"
	"syscall"
	"testing"
	"time"
)

func TestServeCancel(t *testing.T) {
	errCh := make(chan error)
	shutdownTime := time.Millisecond
	srv := func(ctx context.Context, done func()) {
		errCh <- server.Serve(ctx, done, shutdownTime, ":8080",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { time.Sleep(time.Second) }))
	}

	t.Run("call cancel function", func(t *testing.T) {
		ctx, done := context.WithCancel(context.Background())
		go srv(ctx, done)

		done()
		if err := <-errCh; err != nil {
			t.Errorf("error when shutting down via Cancel fn: %v", err)
		}
	})

	t.Run("raise interrupt", func(t *testing.T) {
		ctx, done := context.WithCancel(context.Background())
		go srv(ctx, done)

		time.Sleep(time.Millisecond)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)

		if err := <-errCh; err != nil {
			t.Errorf("error when shutting down via interrupt: %v", err)
		}
	})

	t.Run("raise interrupt while busy", func(t *testing.T) {
		ctx, done := context.WithCancel(context.Background())
		shutdownTime = time.Millisecond
		go srv(ctx, done)

		go http.Get("http://localhost:8080")

		// wait for the request to start getting handled
		time.Sleep(5 * time.Millisecond)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)

		if err := <-errCh; !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("expected context deadline exceeded error when shutting down via interrupt: %v", err)
		}
	})

}
