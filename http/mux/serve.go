package mux

import (
	"net/http"
	"strings"
)

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var allowed []string

	rPath := r.URL.Path

	// https://www.rfc-editor.org/rfc/rfc9110#OPTIONS
	if rPath == "*" && r.Method == http.MethodOptions {
		allowed = append(allowed, http.MethodOptions)
		for _, rt := range m.routes {
			if !contains(allowed, rt.method) {
				allowed = append(allowed, rt.method)
			}
		}
		m.MethodNotAllowed(allowed).ServeHTTP(w, r)
		return
	}

	if rPath == "" || rPath[0] != '/' {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	urlSegments := strings.Split(rPath[1:], "/")

	for _, rt := range m.routes {
		if found := rt.matches(urlSegments); found {
			if rt.method == r.Method {
				r = rt.bindVars(r, urlSegments)

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

	if !contains(allowed, http.MethodOptions) {
		allowed = append(allowed, http.MethodOptions)
	}

	// send 405 with allowed methods
	m.MethodNotAllowed(allowed).ServeHTTP(w, r)
}

func (rt route) matches(urlSegments []string) bool {
	if sn, un := len(rt.segs), len(urlSegments); un < sn || (!rt.wild && un > sn) {
		return false
	}

	varidxs := rt.varidxs

	for i, rSeg := range rt.segs {
		uSeg := urlSegments[i]

		if len(varidxs) > 0 && varidxs[0] == i {
			varidxs = varidxs[1:]

			if re, hasRe := rt.res[rSeg]; hasRe && !re.MatchString(uSeg) {
				return false
			}
			continue
		}

		if rSeg != uSeg {
			return false
		}
	}

	if len(rt.varidxs) == 0 && !rt.wild {
		return true
	}

	return true
}

func (rt *route) bindVars(r *http.Request, urlSegments []string) *http.Request {
	vars := map[string]string{}
	if rt.wild {
		vars["..."] = strings.Join(urlSegments[len(rt.segs):], "/")
	}
	for _, i := range rt.varidxs {
		vars[rt.segs[i]] = urlSegments[i]
	}

	return BindVars(r, vars)
}
