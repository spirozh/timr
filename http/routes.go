package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func Selma(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, "hello selma!!\n")
}

func Long(w http.ResponseWriter, r *http.Request) {
	// Make sure that the writer supports flushing.
	//
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	for i := 0; ; i++ {

		fmt.Fprintf(w, "data: %d\n\n", i)
		flusher.Flush()

		time.Sleep(time.Second)
	}
}

func routes() http.Handler {
	m := mux.NewRouter()
	m.MethodNotAllowedHandler = MethodNotAllowedHandler()
	m.HandleFunc("/selma", Selma).Methods(http.MethodOptions, http.MethodGet, http.MethodOptions)
	m.HandleFunc("/long", Long)

	ApiRoutes(m.PathPrefix("/api").Subrouter())

	printPaths(m)
	return m
}

func MethodNotAllowedHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Header().Add("Allow:", AllowedMethods(r))

		fmt.Fprintf(w, "Method not allowed\n")
		io.WriteString(w, AllowedMethods(r))
	})
}

func ApiRoutes(mux *mux.Router) {
	mux.HandleFunc("/", xxx)
	mux.HandleFunc("/{w}/{z}/", vars)
}

func vars(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	for k, v := range v {
		fmt.Fprintf(w, "%s: %s\n", k, v)
	}
}

func xxx(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "xxx\n")
}

var rootRouter *mux.Router

func AllowedMethods(r *http.Request) string {
	var allowedMethods []string

	rootRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		var m mux.RouteMatch
		if route.Match(r, &m) || m.MatchErr == mux.ErrMethodMismatch {
			methods, err := route.GetMethods()
			if err != nil {
				return err
			}
			allowedMethods = append(allowedMethods, methods...)
		}
		return nil
	})

	return strings.Join(allowedMethods, ", ")
}

func printPaths(r *mux.Router) {
	skip := true
	if skip {
		return
	}

	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})

	if err != nil {
		fmt.Println("Err:", err)
	}
}
