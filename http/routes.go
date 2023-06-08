package http

import (
	"io"
	"net/http"
)

var w string = "aaaa"

func Selma(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "hello selma!!")
}

func routes() http.Handler {
	r := http.NewServeMux()
	r.Handle("/selma", http.HandlerFunc(Selma))

	addRoutes := func(root string, f func(string) http.Handler) {
		h := f(root)
		r.Handle(root, h)
	}
	addRoutes("/api/", ApiRoutes)

	return r
}
