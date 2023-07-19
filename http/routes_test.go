package http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spirozh/timr/test"
)

type selmaTestCase struct {
	method       string
	requestBody  string
	status       int
	responseBody string
}

func Test_Selma(t *testing.T) {

	testCases := []selmaTestCase{
		{http.MethodGet, "", http.StatusOK, "hello selma!!\n"},
		{http.MethodHead, "", http.StatusOK, ""},
		{http.MethodDelete, "", http.StatusMethodNotAllowed, ""},
		{http.MethodPost, "garbage", http.StatusMethodNotAllowed, ""},
		{http.MethodPut, "garbage", http.StatusMethodNotAllowed, ""},
	}

	for _, testCase := range testCases {
		var requestBody io.Reader
		if testCase.requestBody != "" {
			requestBody = strings.NewReader(testCase.requestBody)
		}
		r := httptest.NewRequest(testCase.method, "/selma", requestBody)
		w := httptest.NewRecorder()

		selma(w, r)
		res := w.Result()

		test.Equal(t, testCase.status, res.StatusCode)
		if res.StatusCode == http.StatusMethodNotAllowed {
			test.Equal(t, "GET, HEAD", res.Header["Allow"][0])
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		test.Equal(t, nil, err)
		test.Equal(t, testCase.responseBody, string(data))

		fmt.Printf("%#v\n", res.Header)
	}

}
