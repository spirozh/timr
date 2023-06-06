package http

import (
	"fmt"
	"log"
	"net/http"
)

func Serve() {
	defer fmt.Println("quitting serve")
	http.Handle("/", routes())

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(),
	}

	log.Fatal(srv.ListenAndServe())
}
