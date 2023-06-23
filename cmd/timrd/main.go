package main

import (
	"fmt"
	"time"

	"github.com/spirozh/timr/http"
	"github.com/spirozh/timr/timer"
)

func main() {
	fmt.Println("starting timrd")
	ts := timer.TimerService(time.Now)

	// debugging state
	go func() {
		for {
			fmt.Printf("timers: %v\n\n", ts.List())
			time.Sleep(time.Second)
		}
	}()

	http.Serve(ts)
}
