package mux_test

import (
	"net/http"
	"net/http/httptest"
	"spirozh/timr/http/mux"
	"testing"
)

func mockmiddle(f func()) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			f()
			h.ServeHTTP(w, r)
		})
	}

}
func TestMuxGroup(t *testing.T) {
	var s string

	mm := func(mark string) mux.Middleware {
		return mockmiddle(func() {
			s += mark
		})
	}

	m := mux.New()
	m.Use(mm("a"))

	m1 := m.Group()
	m1.Use(mm("1"))

	m2 := m.Group()
	m2.Use(mm("2"))

	tests := map[string]struct {
		m    *mux.Mux
		want string
	}{
		"base": {m, "a"},
		"one":  {m1, "a1"},
		"two":  {m2, "a2"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			path := "/" + name

			nullHandler := func(w http.ResponseWriter, r *http.Request) {}
			test.m.HandleFunc(path, nullHandler)

			s = ""
			r, _ := http.NewRequest(http.MethodGet, path, nil)
			w := httptest.NewRecorder()
			m.ServeHTTP(w, r)

			if want, got := test.want, s; want != got {
				t.Errorf("s:\nwant: %#v\n got: %#v", want, got)
			}
		})
	}
}
