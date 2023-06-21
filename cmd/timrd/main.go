package main

import (
	"fmt"

	"github.com/spirozh/timr/http"
)

type App struct {
	//htmx *htmx.HTMX
}

func main() {
	fmt.Println("starting timrd")
	http.Serve()
}
