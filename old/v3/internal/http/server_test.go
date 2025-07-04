package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"syscall"
	"testing"
	"time"
)

func TestBadConfig(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())

	s := Server(ctx)

	err := s.Serve()

	netAddrError := &net.AddrError{}
	if !errors.As(err, &netAddrError) {
		t.Fatal(err)
	}
}

func TestServeCancel(t *testing.T) {
	errCh := make(chan error)
	//	shutdownTime := time.Millisecond
	//	responseTime := 10 * time.Millisecond

	srv := func(ctx context.Context, done func()) {
		s := Server(ctx)
		errCh <- s.Serve()
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
		go srv(ctx, done)

		// wait for server to start and then make a request
		time.Sleep(10 * time.Millisecond)
		go http.Get("http://localhost:8080")

		// wait for the request to start getting handled and then send an interrupt
		time.Sleep(5 * time.Millisecond)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)

		if err := <-errCh; !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("expected context deadline exceeded error when shutting down via interrupt: %v", err)
		}
	})
}
