package http

import (
	"fmt"
	"io"
	"net/http"
)

func init() {
	fmt.Println("init api" + w)
}

func fuckThis(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "fuck this\n")
}

func ApiRoutes(root string) http.Handler {
	mux := http.NewServeMux()

	addHandlerFunc := func(path string, h http.HandlerFunc) {
		mux.HandleFunc(root+path, h)
	}

	addHandlerFunc("x/", http.HandlerFunc(fuckThis))

	return mux
}
