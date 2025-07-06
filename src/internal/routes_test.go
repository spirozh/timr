package internal

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	mux := http.NewServeMux()
	a := NewApp()
	a.AddRoutes(mux)

	testCases := map[string]struct {
		path string
		pat  string
	}{
		"index":  {"/", "GET /{$}"},
		"404":    {"/foo", ""},
		"kill":   {"/kill", "GET /kill"},
		"styles": {"/static/styles.css", "GET /static/"},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			r, err := http.NewRequest("GET", "http://localhost"+testCase.path, nil)
			assert.NoError(t, err)
			_, pat := mux.Handler(r)
			assert.Equal(t, testCase.pat, pat)
		})
	}
}
