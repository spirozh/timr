package mux

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var allMethods = []string{
	http.MethodConnect, http.MethodDelete, http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace,
}

func (m *Mux) HandleFunc(path string, handler http.HandlerFunc, methods ...string) {
	m.Handle(path, handler, methods...)
}

func (m *Mux) Handle(path string, handler http.Handler, methods ...string) {
	if path == "" || path[0] != '/' {
		panic(fmt.Sprintf(`invalid path: %#v`, path))
	}
	path = path[1:]

	if contains(methods, http.MethodGet) && !contains(methods, http.MethodHead) {
		methods = append(methods, http.MethodHead)
	}

	if len(methods) == 0 {
		methods = allMethods
	}

	isWild := false
	res := map[string]*regexp.Regexp{}
	varidxs := []int{}

	segs := strings.Split(path, "/")
	for i, seg := range segs {
		switch {
		case len(seg) == 0:
			continue
		case seg == "...":
			if i < len(segs)-1 {
				panic("cannot use wildcard here")
			}
			isWild = true
			segs = segs[:len(segs)-1]
		case seg[0] == ':':
			varname, re, hasRe := strings.Cut(seg[1:], "|")
			for _, i := range varidxs {
				if segs[i] == varname {
					panic("duplicate varname")
				}
			}
			varidxs = append(varidxs, i) // signal that this is a var
			if hasRe {
				res[varname] = regexp.MustCompile(fmt.Sprintf(`^%s$`, re))
			}
			segs[i] = varname
		}
	}

	r := route{
		h:       m.wrap(handler),
		segs:    segs,
		res:     res,
		varidxs: varidxs,
		wild:    isWild,
	}
	for _, method := range methods {
		rCopy := r
		rCopy.method = method
		*m.routes = append(*m.routes, rCopy)
	}
}

func contains(list []string, item string) bool {
	for _, s := range list {
		if item == s {
			return true
		}
	}
	return false
}
