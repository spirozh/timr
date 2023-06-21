package main

import (
	"fmt"
	"time"

	"github.com/spirozh/timr/http"
	"github.com/spirozh/timr/timer"
)

func main() {
	fmt.Println("starting timrd")
	http.Serve(timer.TimerService(time.Now))
}
