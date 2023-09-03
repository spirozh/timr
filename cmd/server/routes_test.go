package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestRoutes(t *testing.T) {
	app := App{
		tokens: map[string]sseSession{},
	}
	ctx, cancel := context.WithCancel(context.Background())
	r := app.Routes(ctx, cancel)

	serveName := func(t *testing.T, path string, headers map[string]string, wantBody string, wantStatus int) {
		t.Helper()
		req, _ := http.NewRequest(http.MethodGet, path, nil)
		for k, v := range headers {
			req.Header.Add(k, v)
		}

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		gotBody := w.Body.String()
		if wantBody != gotBody {
			t.Error("\n(body) want: ", wantBody, "\n        got: ", gotBody)
		}
	}

	t.Run("selma", func(t *testing.T) {
		serveName(t, "/selma", nil, "hello selma", 200)
	})

	t.Run("SSE and trigger", func(t *testing.T) {
		rSSE, _ := http.NewRequest(http.MethodGet, "/SSE", nil)
		wSSE := httptest.NewRecorder()
		go r.ServeHTTP(wSSE, rSSE)

		time.Sleep(time.Millisecond / 10)
		firstData := wSSE.Body.String()

		if !strings.Contains(firstData, "event: token") {
			t.Fatal("no token")
		}
		token := regexp.MustCompile(`data: (.*)\n`).FindStringSubmatch(firstData)[1]

		t.Run("without token", func(t *testing.T) {
			serveName(t, "/trigger", nil, "", http.StatusUnauthorized)

			secondData, _ := strings.CutSuffix(wSSE.Body.String(), firstData)
			if strings.Contains(secondData, "data: triggered") {
				t.Errorf("%q should not contain %q", secondData, "data: triggered")
			}
		})

		t.Run("with token", func(t *testing.T) {
			serveName(t, "/trigger", map[string]string{"Timr-Token": token}, "", http.StatusOK)

			time.Sleep(time.Millisecond / 10)

			secondData, _ := strings.CutSuffix(wSSE.Body.String(), firstData)
			if !strings.Contains(secondData, "data: triggered") {
				t.Errorf("%q should contain %q", secondData, "data: triggered")
			}
		})

		t.Run("shutdown closes SSE", func(t *testing.T) {
			serveName(t, "/shutdown", nil, "", 200)

			time.Sleep(time.Millisecond / 10)
			if ch, exists := app.tokens[token]; exists {
				t.Errorf("token %q should not exist. %v", token, ch)
			}
		})
	})
}
