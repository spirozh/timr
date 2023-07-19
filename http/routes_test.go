package http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/spirozh/timr/html"
	"github.com/spirozh/timr/test"
	"github.com/spirozh/timr/timer"
)

type selmaTestCase struct {
	method       string
	path         string
	requestBody  string
	status       int
	responseBody string
}

func Test_Routes(t *testing.T) {

	indexR, err := html.FS.Open("index.html")
	test.Equal(t, nil, err)
	index, _ := io.ReadAll(indexR)

	testCases := []selmaTestCase{
		// filesystem
		{http.MethodGet, "/", "", http.StatusOK, string(index)},
		{http.MethodPost, "/", "ignored", http.StatusOK, string(index)},  // All methods act like GET
		{http.MethodTrace, "/", "ignored", http.StatusOK, string(index)}, // All methods act like GET
		{http.MethodHead, "/", "", http.StatusOK, ""},                    // Except, of course, HEAD

		{http.MethodGet, "/js", "", http.StatusMovedPermanently, ""},
		{http.MethodGet, "/js/", "", http.StatusOK, "<pre>\n<a href=\"timr.js\">timr.js</a>\n</pre>\n"},

		{http.MethodGet, "/zz", "", http.StatusNotFound, "404 page not found\n"},

		// selma
		{http.MethodGet, "/selma/", "", http.StatusOK, "hello selma!!\n"},
		{http.MethodHead, "/selma/", "", http.StatusOK, ""},
		{http.MethodDelete, "/selma/", "", http.StatusMethodNotAllowed, ""},
		{http.MethodPost, "/selma/", "ignored", http.StatusMethodNotAllowed, ""},
		{http.MethodPut, "/selma/", "ignored", http.StatusMethodNotAllowed, ""},
		{http.MethodPatch, "/selma/", "ignored", http.StatusMethodNotAllowed, ""},

		// redirects to selma
		{http.MethodGet, "/selma", "", http.StatusMovedPermanently, "<a href=\"/selma/\">Moved Permanently</a>.\n\n"},
		{http.MethodPost, "/selma", "ignored", http.StatusMovedPermanently, ""},
	}

	ts := timer.TimerService(func() time.Time { return time.Time{} })
	routes := routes(ts)

	for i, testCase := range testCases {
		fmt.Printf("testcase: %d) %#v\n", i, testCase)
		var requestBody io.Reader
		if testCase.requestBody != "" {
			requestBody = strings.NewReader(testCase.requestBody)
		}
		r := httptest.NewRequest(testCase.method, testCase.path, requestBody)
		w := httptest.NewRecorder()

		routes.ServeHTTP(w, r)
		res := w.Result()

		test.Equal(t, testCase.status, res.StatusCode)
		if res.StatusCode == http.StatusMethodNotAllowed {
			test.Equal(t, "GET, HEAD", res.Header["Allow"][0])
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		test.Equal(t, nil, err)
		test.Equal(t, testCase.responseBody, string(data))
	}
}
