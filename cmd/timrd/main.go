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
	sub := ts.Subscribe(func(_ timr.TimrEventType, _ string, _ timr.Timer) {

		var b bytes.Buffer
		list := ts.List()
		fmt.Fprintf(&b, "timers: %v\n", list)
		for _, name := range list {
			t, err := ts.Get(name)
			if err != nil {
				fmt.Fprintf(&b, " ! %v: %v\n", name, err)
				continue
			}
			fmt.Fprintf(&b, " * %v: %#v\n", name, t.State())
		}
		fmt.Fprintln(&b)

		timr.INFO("current timerServer state:\n", b.String())
	})

	http.Serve(ts)

	ts.Unsubscribe(sub)
}
