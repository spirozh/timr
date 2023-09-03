package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"spirozh/timr/http/mux"
)

type App struct {
	tokens map[string]chan string
}

func (app *App) Routes(ctx context.Context, cancel func()) http.Handler {
	m := mux.New()

	m.HandleFunc("/selma", Selma, http.MethodGet)
	m.Handle("/shutdown", Shutdown(cancel), http.MethodGet)
	m.Handle("/SSE", app.SSE(ctx), http.MethodGet)

	m.Use(TimrToken())
	m.Handle("/trigger", app.Trigger())

	return m
}

func TimrToken() mux.Middleware {
	return func(h http.Handler) http.Handler {
		return h
	}
}

func Selma(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello selma"))
}

func Shutdown(cancel func()) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		cancel()
	}
}

func (app *App) Trigger() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Timr-Token")
		if ch, ok := app.tokens[token]; ok {
			go func() { ch <- "triggered" }()
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

func init() {
	// Panic if a cryptographically secure PRNG is not available.
	buf := make([]byte, 1)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: Read() failed with %#v", err))
	}
}

func RandomToken(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (app *App) SSE(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, _ := w.(http.Flusher)

		// write headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// generate token
		tok := RandomToken(5)
		// send token
		fmt.Fprintf(w, "event: token\ndata: %s\n\n", tok)
		flusher.Flush()

		// save channel with token
		app.tokens[tok] = make(chan string)
		defer delete(app.tokens, tok)

		for {
			select {
			case data := <-app.tokens[tok]:
				fmt.Fprintf(w, "data: %v\n\n", data)
				flusher.Flush()
			case <-ctx.Done():
				return
			case <-r.Context().Done():
				return
			}
		}

	}

}
