package main

import (
	"net/http"
	"spirozh/timr/internal/http/mux"
)

func main() {
	app := App()
	app.run()
}

func routes(cancel func()) http.Handler {
	m := mux.New()

	return m
}
