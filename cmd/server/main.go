package main

import (
	"context"
	"fmt"
	"spirozh/timr/http"
	"spirozh/timr/timer"
	"time"
)

func main() {
	t := timer.New()
	t.Start(time.Now())

	app := App{
		tokens: map[string]chan string{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	http.Serve(ctx, cancel, app.Routes(ctx, cancel))

	fmt.Println(t.Elapsed(time.Now()))
}
