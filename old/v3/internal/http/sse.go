package http

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"sync"
)

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
	Event string
	Data  string
}

func (e SSEEvent) Write(w io.Writer) {
	if e.Event != "" {
		fmt.Fprintf(w, "event: %s\n", e.Event)
	}
	fmt.Fprintf(w, "data: %s\n\n", e.Data)

	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

type SSETokenMap struct {
	mu       sync.RWMutex
	channels map[string]chan SSEEvent
}

<<<<<<<< HEAD:src/internal/server/sse.go
func (t SSETokenMap) SSE(serverCtx context.Context) http.HandlerFunc {
========
func (t *SSETokenMap) newToken() string {
	t.mu.Lock()
	defer t.mu.Unlock()

	var token string
	for {
		token = RandomToken(15)
		if _, used := t.channels[token]; !used {
			break
		}
	}
	t.channels[token] = make(chan SSEEvent)
	return token
}

func (t *SSETokenMap) getChannel(token string) chan SSEEvent {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.channels[token]
}

func (t *SSETokenMap) discardToken(token string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	close(t.getChannel(token))
	delete(t.channels, token)
}

func (t *SSETokenMap) SSE(serverCtx context.Context, sseChannels map[string]chan SSEEvent) http.HandlerFunc {
>>>>>>>> 21b71bc (start again):old/v3/internal/http/sse.go
	return func(w http.ResponseWriter, r *http.Request) {
		// write headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		tok := t.newToken()

		// send token
		SSEEvent{Event: "SSE-Token", Data: tok}.Write(w)

		ch := t.getChannel(tok)
		defer t.discardToken(tok)

		requestCtx := r.Context()
		for {
			select {
			case event := <-ch:
				event.Write(w)
				continue
			case <-serverCtx.Done():
			case <-requestCtx.Done():
			}
			break
		}
	}
}

type sessionKeyType struct{}

var sessionKey sessionKeyType

<<<<<<<< HEAD:src/internal/server/sse.go
func (t SSETokenMap) RequireSSEToken(h http.Handler) http.Handler {
========
func (t *SSETokenMap) RequireSSEToken(h http.Handler) http.Handler {
>>>>>>>> 21b71bc (start again):old/v3/internal/http/sse.go
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-SSE-Token")
		if ch := t.getChannel(token); ch != nil {
			r = r.WithContext(context.WithValue(r.Context(), sessionKey, ch))
			h.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
	})
}

func SSEChan(r *http.Request) chan<- SSEEvent {
	ch := r.Context().Value(sessionKey)
	if eChan, ok := ch.(chan SSEEvent); ok {
		return eChan
	}
	return nil
}
