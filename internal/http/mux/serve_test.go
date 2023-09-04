package mux_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"spirozh/timr/http/mux"
	"testing"
)

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
				{http.MethodHead, "/", http.StatusOK, nil, ""}, // registering GET should register HEAD
				{http.MethodOptions, "/", http.StatusNoContent, nil, "GET, HEAD, OPTIONS"},
				{http.MethodPost, "/", http.StatusMethodNotAllowed, nil, "GET, HEAD, OPTIONS"},
				{http.MethodGet, "", http.StatusBadRequest, nil, ""},
			},
		},
		{
			"/:k", []string{http.MethodGet},
			[]subtest{
				{http.MethodGet, "/v", http.StatusOK, map[string]string{"k": "v"}, ""},
				{http.MethodGet, "/", http.StatusOK, map[string]string{"k": ""}, ""},
				{http.MethodGet, "http://host//", http.StatusNotFound, nil, ""}, // to get '//' as path, urlParse needs more..
			},
		},
		{
			`/:k|\d`, []string{http.MethodGet},
			[]subtest{
				{http.MethodGet, "/1", http.StatusOK, map[string]string{"k": "1"}, ""},
				{http.MethodGet, "/v", http.StatusNotFound, nil, ""},
				{http.MethodGet, "/", http.StatusNotFound, nil, ""},
				{http.MethodGet, "http://host//", http.StatusNotFound, nil, ""},
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
						t.Errorf("status\n want: %#v\n  got: %#v", sub.status, status)
					}

					allow := w.Result().Header.Get("Allow")
					if allow != sub.allow {
						t.Errorf("allow\n want: %#v\n  got: %#v", sub.allow, allow)
					}

					if !paramsDone {
						t.Error("params were xx")
					}
				})
			}
		})
	}
}
