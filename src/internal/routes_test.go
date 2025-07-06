package internal

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStaticRouteVsIndex(t *testing.T) {
	mux := http.NewServeMux()
	a := NewApp()
	a.AddRoutes(mux)

	testCases := map[string]struct {
		method string
		path   string
		pat    string
	}{
		"index":  {"GET", "/", "GET /"},
		"kill":   {"GET", "/kill", "GET /kill"},
		"styles": {"GET", "/static/styles.css", "GET /static/"},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			r, err := http.NewRequest(testCase.method, "http://localhost"+testCase.path, nil)
			assert.NoError(t, err)
			_, pat := mux.Handler(r)
			assert.Equal(t, testCase.pat, pat)
		})
	}
}
