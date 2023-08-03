package mux

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

type mux struct {
	routes     []route
	handler404 http.Handler
	handler405 http.Handler
}

type route struct {
	h       http.Handler
	method  string
	segs    []string
	wild    bool
	varidxs []int
	res     map[string]*regexp.Regexp
}

func New() http.Handler {
	return &mux{}
}

var allMethods = []string{http.MethodConnect, http.MethodDelete, http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace}

func (m *mux) Handle(path string, handler http.Handler, methods ...string) {
	if len(methods) == 0 {
		methods = allMethods
	}

	var (
		isWild  bool
		res     map[string]*regexp.Regexp
		varidxs []int
	)

	segs := strings.Split(path, "/")
	for i, seg := range segs {
		switch {
		case seg[0] == ':':
		case seg == "...":
			if i < len(segs)-1 {
				panic("cannot use wildcard here")
			}
			isWild = true

		case seg[0] == ':':
			varidxs = append(varidxs, i) // signal that this is a var
			// check for regex
			varname, re, hasRe := strings.Cut(seg[1:], "|")
			if hasRe {
				compRe, err := regexp.Compile(re)
				if err != nil {
					panic("or something")
				}
				if res == nil {
					res = map[string]*regexp.Regexp{}
				}
				res[varname] = compRe
			}
			segs[i] = varname
		}
	}

	r0 := route{
		h:       handler,
		segs:    segs,
		res:     res,
		varidxs: varidxs,
		wild:    isWild,
	}
	for _, method := range methods {
		r1 := r0
		r1.method = method
		m.routes = append(m.routes, r1)
	}

}

func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var allowed []string
	for _, rt := range m.routes {
		if vars, found := rt.matches(r.URL.Path); found {
			if rt.method == r.Method {
				// bind vars
				r = r.WithContext(context.WithValue(r.Context(), varkey, vars))

				rt.h.ServeHTTP(w, r)
				return
			}
			allowed = append(allowed, r.Method)
		}
	}

	if allowed == nil {
		// send 404
		w.WriteHeader(http.StatusNotFound)
	}

	// send 405 with allow methods
	w.Header().Add("Allow", strings.Join(allowed, ","))
	w.WriteHeader(http.StatusMethodNotAllowed)

}

type varkeytype struct{}

var varkey varkeytype

func Var(r *http.Request, key string) (val string, ok bool) {
	varsAny := r.Context().Value(varkey)
	if varsAny == nil {
		return "", false
	}

	vars, ok := varsAny.(map[string]string)
	if !ok {
		return "", false
	}

	return vars[key], true

}

func (rt route) matches(path string) (map[string]string, bool) {
	urlSegments := strings.Split(path, "/")

	if sn, un := len(rt.segs), len(urlSegments); un < sn || (!rt.wild && un > sn) {
		return nil, false
	}

	varidxs := rt.varidxs

	for i, rSeg := range rt.segs {
		uSeg := urlSegments[i]

		if len(varidxs) > 0 && varidxs[0] == i {
			varidxs = varidxs[1:]

			if re, hasRe := rt.res[rSeg]; hasRe {
				if !re.MatchString(uSeg) {
					return nil, false
				}
			}

			continue
		}

		if rSeg != uSeg {
			return nil, false
		}
	}

	if len(rt.varidxs) == 0 {
		return nil, true
	}

	vars := map[string]string{}
	if rt.wild {
		vars["..."] = strings.Join(urlSegments[len(rt.segs):], "/")
	}
	for _, i := range rt.varidxs {
		vars[rt.segs[i]] = urlSegments[i]
	}
	return vars, true
}
