package mux_test

import (
	"net/http"
	"net/http/httptest"
	"spirozh/timr/internal/mux"
	"testing"
)

func TestMiddleWare(t *testing.T) {
	var s string

	m1 := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s += "1"
			h.ServeHTTP(w, r)
		})
	}
	m2 := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s += "2"
			h.ServeHTTP(w, r)
		})
	}

	h := func(w http.ResponseWriter, r *http.Request) { s += "h" }

	tests := map[string]struct {
		mws  []mux.Middleware
		want string
	}{
		"zero": {[]mux.Middleware{}, "h"},
		"one":  {[]mux.Middleware{m1}, "1h"},
		"two":  {[]mux.Middleware{m1, m2}, "12h"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := mux.New()
			m.Use(test.mws...)
			m.HandleFunc("/", h)

			s = ""
			r, _ := http.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			m.ServeHTTP(w, r)

			if want, got := test.want, s; want != got {
				t.Errorf("s:\nwant: %#v\n got: %#v", want, got)
			}
		})
	}
}
