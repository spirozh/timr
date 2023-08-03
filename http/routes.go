package http

import (
	"net/http"
	"spirozh/timr/http/mux"
)

func routes(shutdown chan<- struct{}) http.Handler {
	m := mux.New()

	m.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		close(shutdown)
	}, http.MethodGet)

	return m
}
