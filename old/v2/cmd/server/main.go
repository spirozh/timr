package main

import (
	"context"
	"fmt"
	"spirozh/timr/old/v2/internal/http"
	"spirozh/timr/old/v2/internal/timer"
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
