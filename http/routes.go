package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Selma(r *mux.Router) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "hello selma!!")

		r.Walk(func(route *mux.Route, router *mux.Router, ancestor []*mux.Route) error {
			return nil
		})
	}
}

func routes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/health", HealthCheckHandler)
	r.HandleFunc("/", Selma(r))
	r.Use()
	return r
}
