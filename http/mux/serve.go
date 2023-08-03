package mux

import (
	"context"
	"net/http"
	"strings"
)

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
