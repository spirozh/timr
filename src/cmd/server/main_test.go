package main_test

import (
	"errors"
	"io"
	"log"
	"net/http"
	"syscall"
	"testing"

	"github.com/spirozh/timr/internal"
)

func TestMain(m *testing.M) {
	a := internal.NewApp()
	a.Run()
	defer func() {
		if err := a.Close(); err != nil {
			log.Fatal("closing app: ", err)
		}
	}()

	m.Run()
}

func TestRoutes(t *testing.T) {
	isServerAlive(t, true)

	expected := "Hello, Selma!\n"
	hasBody(t, "http://localhost:8080/selma", expected)

	expected = "Goodbye\n"
	hasBody(t, "http://localhost:8080/kill", expected)
	isServerAlive(t, false)
}

func isServerAlive(t *testing.T, alive bool) {
	_, err := http.Get("http://localhost:8080/")
	if alive && err != nil {
		t.Error(err)
	}

	if !alive && !errors.Is(err, syscall.ECONNREFUSED) {
		t.Error("expected syscall.ECONNREFUSED, got: ", err)
	}
}

func hasBody(t *testing.T, url string, expected string) {
	r, err := http.Get(url)
	if err != nil {
		t.Error(err)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Error(err)
		return
	}
	if string(b) != expected {
		err := errors.New("expected '" + expected + "', got: " + string(b))
		t.Error(err)
		return
	}
}
