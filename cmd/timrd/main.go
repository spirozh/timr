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
	var i int
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			list := ts.List()
			fmt.Printf("%05d) timers: %v\n", i, list)
			i++

			for _, name := range list {
				t, err := ts.Get(name)
				if err != nil {
					fmt.Printf(" ! %v: %v\n", name, err)
					continue
				}
				fmt.Printf(" * %v: %v\n", name, t.State())
			}

			fmt.Println()
			<-ticker.C
		}
	}()

	http.Serve(ts)
}
