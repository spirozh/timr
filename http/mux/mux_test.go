package mux_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestMuxMatch(t *testing.T) {

	type subtest struct {
		method  string
		request string

		status int
		params map[string]string
		allow  string
	}
	tests := []struct {
		path    string
		methods []string

		subtests []subtest
	}{
		{
			"/", []string{http.MethodGet},
			[]subtest{
				{http.MethodGet, "/", http.StatusOK, nil, ""},
				// {http.MethodHead, "/", http.StatusOK, nil, ""},    // registering GET should register HEAD
				// {http.MethodOptions, "/", http.StatusOK, nil, ""}, // without registered OPTIONS, automatisch
				{http.MethodPost, "/", http.StatusMethodNotAllowed, nil, "GET"},
				{http.MethodGet, "", http.StatusNotFound, nil, ""},
			},
		},
		{
			"/:k", []string{http.MethodGet},
			[]subtest{
				{http.MethodGet, "/v", http.StatusOK, map[string]string{"k": "v"}, ""},
				{http.MethodGet, "/", http.StatusOK, map[string]string{"k": ""}, ""},
				{http.MethodGet, "//", http.StatusNotFound, nil, ""},
			},
		},
		{
			`/:k|\d`, []string{http.MethodGet},
			[]subtest{
				{http.MethodGet, "/1", http.StatusOK, map[string]string{"k": "1"}, ""},
				{http.MethodGet, "/v", http.StatusNotFound, nil, ""},
				{http.MethodGet, "/", http.StatusNotFound, nil, ""},
				{http.MethodGet, "//", http.StatusNotFound, nil, ""},
			},
		},
	}
	for _, test := range tests {
		name := fmt.Sprintf("%s%s", test.methods, test.path)
		t.Run(name, func(t *testing.T) {
			var (
				params     map[string]string
				paramsDone bool
				subT       *testing.T
			)

			m := mux.New()
			m.HandleFunc(test.path, func(w http.ResponseWriter, r *http.Request) {
				if params == nil {
					return
				}

				for k, v := range params {
					if val, found := mux.Var(r, k); found {
						if v != val {
							subT.Errorf(`param["%s"]="%s", expected "%s"`, k, val, v)
						}
						continue
					}
					subT.Errorf(`param["%s"] not found`, k)
				}
				paramsDone = true
			}, test.methods...)

			for _, sub := range test.subtests {
				name := fmt.Sprintf("%s %s", sub.method, sub.request)
				t.Run(name, func(t *testing.T) {
					params, subT = sub.params, t
					paramsDone = (params == nil)
					r, _ := http.NewRequest(sub.method, sub.request, nil)
					w := httptest.NewRecorder()
					m.ServeHTTP(w, r)

					status := w.Result().StatusCode
					if status != sub.status {
						t.Error("status\n want:", sub.status, "\n  got:", status)
					}

					allow := w.Result().Header.Get("Allow")
					if allow != sub.allow {
						t.Error("allow\n want:", sub.allow, "\n  got:", allow)
					}

					if !paramsDone {
						t.Error("params were xx")
					}
				})
			}
		})
	}
}
