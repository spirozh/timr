package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
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

type SSETokenMap map[string]chan SSEEvent

func (t SSETokenMap) SSE(serverCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// write headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// generate not currently used token (should also check for 'recently used'?)
		var tok string
		for {
			tok = RandomToken(15)
			if _, alreadyExists := t[tok]; !alreadyExists {
				break
			}
		}

		// send token
		SSEEvent{Event: "SSE-Token", Data: tok}.Write(w)

		ch := make(chan SSEEvent)
		defer close(ch)

		// save channel with token for for sending events to this connections
		t[tok] = ch
		defer delete(t, tok)

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

func (t SSETokenMap) RequireSSEToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-SSE-Token")
		ch, ok := t[token]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), sessionKey, ch))
		h.ServeHTTP(w, r)
	})
}

func SSEChan(r *http.Request) chan<- SSEEvent {
	ch := r.Context().Value(sessionKey)
	if eChan, ok := ch.(chan SSEEvent); ok {
		return eChan
	}
	return nil
}
