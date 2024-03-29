package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spirozh/timr"
)

func Serve(ts timr.TimerService) {
	defer fmt.Println("quitting serve")

	sseDone := make(chan any)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(ts, sseDone),
	}

	go func() {
		timr.INFO("listening on", srv.Addr)

		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	waitForShutdown(srv, sseDone)
}

func waitForShutdown(srv *http.Server, sseDone chan any) {
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Shutdown SSE connections
	close(sseDone)

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
