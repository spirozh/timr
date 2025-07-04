package http_test

import (
	"bytes"
	"spirozh/timr/internal/http"
	"testing"
)

type flushableBuffer struct {
	bytes.Buffer
	flushCount int
}

func (b *flushableBuffer) Flush() {
	b.flushCount++
}

func TestSSEEventWrite(t *testing.T) {
	for name, test := range map[string]struct {
		event  http.SSEEvent
		output string
	}{
		"zero event":     {http.SSEEvent{}, "data: \n\n"},
		"data x":         {http.SSEEvent{Data: "x"}, "data: x\n\n"},
		"data x event y": {http.SSEEvent{Data: "x", Event: "y"}, "event: y\ndata: x\n\n"},
	} {
		t.Run(name, func(t *testing.T) {
			var b flushableBuffer
			test.event.Write(&b)
			if b.String() != test.output {
				t.Errorf("b.String()\nwant: %#v\n got: %#v", test.output, b.String())
			}
			if b.flushCount != 1 {
				t.Errorf("b.flushCount\nwant: 1\n got: %d", b.flushCount)
			}
		})
	}
}
