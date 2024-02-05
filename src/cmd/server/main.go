package main

import (
	"context"
	"fmt"
	"spirozh/timr/internal/mux"
	"spirozh/timr/internal/server"
	"time"
)

func main() {
	ctx, done := context.WithCancel(context.Background())

	h := mux.New()

	fmt.Println("helloooo")
	server.Serve(ctx, done, time.Second, "localhost:8080", h)
	fmt.Println("goodbyyyeeeee")
}
