package mux_test

import (
	"net/http"
	"spirozh/timr/http/mux"
	"testing"
)

func TestMuxHandleGood(t *testing.T) {
	var dummy http.HandlerFunc
	m := mux.New()
	paths := map[string]string{
		"root":           "/",
		"one static seg": "/foo",
		"one var":        "/:foo",
		"regex":          "/:foo|[0-9]+",
		"wildcard":       "/...",
		"two vars":       "/:foo/:bar",
	}
	for testName, path := range paths {
		t.Run(testName, func(t *testing.T) {
			m.Handle(path, dummy, http.MethodGet)
		})
	}
}

func TestMuxHandleBad(t *testing.T) {
	var dummy http.HandlerFunc
	m := mux.New()
	paths := map[string]string{
		"empty":              "",
		"no leading slash":   "foo",
		"duplicate var":      "/:foo/:foo",
		"wildcard in middle": "/.../foo",
		"bad re":             "/:foo|[",
	}
	for name, path := range paths {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Error("no panic")
				}
			}()
			m.Handle(path, dummy, http.MethodGet)
		})
	}
}
