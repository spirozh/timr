package mux

import "net/http"

type Middleware func(http.Handler) http.Handler

func (m *Mux) Use(mws ...Middleware) {
	m.mws = append(m.mws, mws...)
}

func (m *Mux) wrap(h http.Handler) http.Handler {
	for i := len(m.mws) - 1; i >= 0; i-- {
		h = m.mws[i](h)
	}
	return h
}
