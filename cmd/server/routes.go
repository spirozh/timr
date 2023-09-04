package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"spirozh/timr/internal/http/mux"
)

type sessionKeyType struct{}

var sessionKey sessionKeyType

type sseSession struct {
	ch chan SSEEvent
}

type App struct {
	tokens map[string]sseSession
}

func (app *App) Routes(ctx context.Context, cancel func()) http.Handler {
	m := mux.New()
	m.Use(NoPanic)

	m.HandleFunc("/selma", Selma, http.MethodGet)
	m.Handle("/shutdown", Shutdown(cancel), http.MethodGet)
	m.Handle("/SSE", app.SSE(ctx), http.MethodGet)

	m.Use(app.TimrToken)
	m.Handle("/trigger", app.Trigger())

	return m
}

func NoPanic(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			panicked := recover()
			if panicked == nil {
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "\n\n%#v\n\n", panicked)
		}()
		h.ServeHTTP(w, r)
	})
}
func (app *App) TimrToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Timr-Token")
		ss, ok := app.tokens[token]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), sessionKey, ss))
		h.ServeHTTP(w, r)
	})
}

func session(r *http.Request) sseSession {
	ss := r.Context().Value(sessionKey)
	return ss.(sseSession)
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
		go func() { session(r).ch <- SSEEvent{data: "triggered"} }()
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

type SSEEvent struct {
	event string
	data  string
}

func (e SSEEvent) Write(w io.Writer) {
	flusher, _ := w.(http.Flusher)

	if e.event != "" {
		fmt.Fprintf(w, "event: %s\n", e.event)
	}
	fmt.Fprintf(w, "data: %s\n", e.data)
	flusher.Flush()
}

func (app *App) SSE(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// write headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// generate token
		tok := RandomToken(5)
		// send token
		SSEEvent{event: "token", data: tok}.Write(w)

		// save channel with token
		app.tokens[tok] = sseSession{make(chan SSEEvent)}
		defer delete(app.tokens, tok)

		for {
			select {
			case event := <-app.tokens[tok].ch:
				event.Write(w)
			case <-ctx.Done():
				return
			case <-r.Context().Done():
				return
			}
		}

	}

}
