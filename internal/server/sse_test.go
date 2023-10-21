package server_test

import (
	"bytes"
	"spirozh/timr/internal/server"
	"testing"
)

func TestSSEEventWrite(t *testing.T) {
	for name, test := range map[string]struct {
		event  server.SSEEvent
		output string
	}{
		"zero event":     {server.SSEEvent{}, "data: \n\n"},
		"data x":         {server.SSEEvent{Data: "x"}, "data: x\n\n"},
		"data x event y": {server.SSEEvent{Data: "x", Event: "y"}, "event: y\ndata: x\n\n"},
	} {
		t.Run(name, func(t *testing.T) {
			var b bytes.Buffer
			test.event.Write(&b)
			if b.String() != test.output {
				t.Errorf("\nwant: %#v\n got: %#v", test.output, b.String())
			}
		})
	}
}
