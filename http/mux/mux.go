package mux

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

type Mux struct {
	routes     []route
	NotFound   http.Handler
	handler405 func(allowed []string) http.Handler
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
		handler405: func(allowed []string) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Allow", strings.Join(allowed, ","))
				w.WriteHeader(http.StatusMethodNotAllowed)
			})
		},
		NotFound: http.NotFoundHandler(),
	}
}

var allMethods = []string{
	http.MethodConnect, http.MethodDelete, http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace,
}

func (m *Mux) HandleFunc(path string, handler http.HandlerFunc, methods ...string) {
	m.Handle(path, handler, methods...)
}

func (m *Mux) Handle(path string, handler http.Handler, methods ...string) {
	if path == "" {
		panic("path empty")
	}
	if path[0] != '/' {
		panic("path must start with '/'")
	}

	if len(methods) == 0 {
		methods = allMethods
	}

	var (
		isWild  bool
		res     map[string]*regexp.Regexp
		varidxs []int
	)

	//for seg, path, found := strings.Cut(path, "/");

	path = path[1:]

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
				compRe, err := regexp.Compile("^" + re + "$")
				if err != nil {
					panic("bad regexp")
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

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var allowed []string
	for _, rt := range m.routes {
		if vars, found := rt.matches(r.URL.Path); found {
			if rt.method == r.Method {
				// bind vars
				r = r.WithContext(context.WithValue(r.Context(), varkey, vars))

				rt.h.ServeHTTP(w, r)
				return
			}
			allowed = append(allowed, rt.method)
		}
	}

	if allowed == nil {
		m.NotFound.ServeHTTP(w, r)
		return
	}

	// send 405 with allow methods
	m.handler405(allowed).ServeHTTP(w, r)
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

	val, found := vars[key]
	return val, found

}

func (rt route) matches(path string) (map[string]string, bool) {
	if len(path) == 0 || path[0] != '/' {
		return nil, false
	}
	urlSegments := strings.Split(path[1:], "/")

	if sn, un := len(rt.segs), len(urlSegments); un < sn || (!rt.wild && un > sn) {
		return nil, false
	}

	varidxs := rt.varidxs

	for i, rSeg := range rt.segs {
		uSeg := urlSegments[i]

		if len(varidxs) > 0 && varidxs[0] == i {
			varidxs = varidxs[1:]

			if re, hasRe := rt.res[rSeg]; hasRe && !re.MatchString(uSeg) {
				return nil, false
			}
			continue
		}

		if rSeg != uSeg {
			return nil, false
		}
	}

	if len(rt.varidxs) == 0 && !rt.wild {
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
