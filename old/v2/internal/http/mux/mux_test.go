package mux_test

import (
	"spirozh/timr/internal/http/mux"
	"testing"
)

func TestNew(t *testing.T) {
	m := mux.New()
	if m == nil {
		t.Error("New() returns nil")
	}
}
