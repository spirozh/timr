package main

import (
	"context"
	"fmt"
	"spirozh/timr/internal/http"
	"spirozh/timr/internal/timer"
	"time"
)

func main() {
	t := timer.New()
	t.Start(time.Now())

	app := App{
		tokens: map[string]sseSession{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	http.Serve(ctx, cancel, app.Routes(ctx, cancel))

	fmt.Println(t.Elapsed(time.Now()))
}
