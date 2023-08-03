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

func Serve() {
	shutdown := make(chan struct{})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(shutdown),
	}

	var mu sync.Mutex
	mu.Lock()
	go func() {
		defer mu.Unlock()
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	waitForShutdown(srv, shutdown)
	mu.Lock()
}

func waitForShutdown(srv *http.Server, shutdown chan struct{}) {
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal or shutdown is closed
	select {
	case <-c:
	case <-shutdown:
	}

	closeIfOpen(shutdown)

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline.
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Println("server Shutdown error: ", err)
	}
}

func closeIfOpen(ch chan struct{}) {
	select {
	case <-ch:
	default:
		close(ch)
	}
}
