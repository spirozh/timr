package main

import (
	"context"
	"net/http"
	"spirozh/timr/http/mux"
)

func Routes(ctx context.Context, cancel func()) http.Handler {
	m := mux.New()

	m.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		cancel()
	}, http.MethodGet)

	return m
}
