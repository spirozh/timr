package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/spirozh/timr"
	"github.com/spirozh/timr/http"
	"github.com/spirozh/timr/timer"
)

func main() {
	timr.INFO("starting timrd")
	ts := timer.TimerService(time.Now)

	// debugging state
	sub := ts.Subscribe(func(_ timr.TimrEventType, _ int, _ string, _ timr.Timer) { logState(ts) })
	defer ts.Unsubscribe(sub)

	http.Serve(ts)
}

func logState(ts timr.TimerService) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "timers:\n")
	ts.ForAll(func(id int, name string, timer timr.TimerState) {
		fmt.Fprintf(&b, " * %d) '%s': %v\n", id, name, timer)
	})
	fmt.Fprintln(&b)

	timr.INFO("current timerServer state:\n", b.String())
}
