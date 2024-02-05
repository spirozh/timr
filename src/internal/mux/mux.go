// Package mux implements an HTTP router.
//
// Features:
// * automatic Allow header generation for 405
// * automatic HEAD handler
// * custom 404 and 405 handlers
// * paths with named placeholders (conflicts resolved by first defined) and optional methods
// * middleware
//
// TODO: groups, automatic OPTION handler, documentation
package mux

import (
	"net/http"
	"regexp"
	"strings"
)

type Mux struct {
	routes           *[]route
	mws              []Middleware
	NotFound         http.Handler
	MethodNotAllowed func(allowed []string) http.Handler
}

type route struct {
	h       http.Handler
	method  string
	segs    []string
	wild    bool
	varidxs []int
	res     map[string]*regexp.Regexp
}

func New() *Mux {
	return &Mux{
		routes: &[]route{},
		mws:    []Middleware{},
		MethodNotAllowed: func(allowed []string) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Allow", strings.Join(allowed, ", "))
				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusMethodNotAllowed)
			})
		},
		NotFound: http.NotFoundHandler(),
	}
}
