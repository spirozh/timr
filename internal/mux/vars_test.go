package mux_test

import (
	"net/http"
	"spirozh/timr/internal/mux"
	"testing"
)

func TestSetVar(t *testing.T) {
	r, err := http.NewRequest("", "/", nil)
	if err != nil {
		t.Error("problem making request: ", err)
	}

	k, v := "foo", "bar"
	newr := mux.BindVar(r, k, v)
	if newr == r {
		t.Fatal("SetVar should make a new request the first time")
	}

	got, found := mux.Var(newr, k)
	if !found {
		t.Errorf("couldn't find k: %#v", k)
	}
	if v != got {
		t.Errorf("mux.Var(r, %#v):\nwant: %#v\n got: %#v", k, v, got)
	}

	nextr := mux.BindVar(newr, k, v)
	if newr != nextr {
		t.Fatal("SetVar should not make a new request the second time")
	}
}
