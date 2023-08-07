package mux

import (
	"context"
	"net/http"
)

type varkeytype struct{}

var varkey varkeytype

func BindVars(r *http.Request, vars map[string]string) *http.Request {
	varsAny := r.Context().Value(varkey)
	if varsAny == nil {
		varsAny = map[string]string{}

		ctx := context.WithValue(r.Context(), varkey, varsAny)
		r = r.WithContext(ctx)
	}

	ctxVars := varsAny.(map[string]string)

	for k, v := range vars {
		ctxVars[k] = v
	}

	return r
}

func BindVar(r *http.Request, key string, val string) *http.Request {
	return BindVars(r, map[string]string{key: val})
}

func Var(r *http.Request, key string) (val string, ok bool) {
	varsAny := r.Context().Value(varkey)
	if varsAny == nil {
		return "", false
	}

	vars, ok := varsAny.(map[string]string)
	if !ok {
		panic("bad value in context")
	}

	val, found := vars[key]
	return val, found
}
