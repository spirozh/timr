package mux

import (
	"context"
	"net/http"
)

type varkeytype struct{}

var varkey varkeytype

func SetVar(r *http.Request, key string, val string) *http.Request {
	varsAny := r.Context().Value(varkey)
	if varsAny == nil {
		varsAny = map[string]string{}
		r = r.WithContext(context.WithValue(r.Context(), varkey, varsAny))
	}

	vars := varsAny.(map[string]string)
	vars[key] = val

	return r
}
func Var(r *http.Request, key string) (val string, ok bool) {
	varsAny := r.Context().Value(varkey)
	if varsAny == nil {
		return "", false
	}

	vars, ok := varsAny.(map[string]string)
	if !ok {
		return "", false
	}

	val, found := vars[key]
	return val, found
}
