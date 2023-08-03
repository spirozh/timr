package main

import (
	"fmt"
	"spirozh/timr/timer"
	"time"
)

func main() {

	t := timer.New()
	t.Start(time.Now())
	fmt.Println(t.Elapsed(time.Now()))
}
