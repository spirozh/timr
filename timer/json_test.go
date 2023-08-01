package timer

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func Test_JSONMarshal(t *testing.T) {
	type test struct {
		name  string
		timer timer
		json  string
	}

	tests := []test{
		{
			name:  "zero timer",
			timer: New(),
			json:  `{"name":"","duration":0,"started":"","elapsed":0}`,
		},
		{
			name:  "zero started",
			timer: New(Name("foo"), Duration(time.Second), Started(&time.Time{}), Elapsed(time.Minute)),
			json:  `{"name":"foo","duration":1000,"started":"0001-01-01T00:00:00.0000","elapsed":60000}`,
		},
		{
			name:  "nil started",
			timer: New(Name("bar"), Duration(time.Minute), Started(nil), Elapsed(time.Second)),
			json:  `{"name":"bar","duration":60000,"started":"","elapsed":1000}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := json.Marshal(test.timer)
			if err != nil {
				t.Errorf("err\nexpected: nil\n  actual: %v", err)
			}
			if string(b) != test.json {
				t.Errorf("b\nexpected: %v\n  actual: %v", test.json, string(b))
			}
		})
	}
}

func Test_JSONUnmarshal(t *testing.T) {
	type test struct {
		name         string
		json         string
		timer        timer
		unmarshalErr error
	}

	tests := []test{
		{
			name:  "started zero",
			json:  `{"name":"foo","duration":1000,"started":"0001-01-01T00:00:00.0000","elapsed":60000}`,
			timer: New(Name("foo"), Duration(time.Second), Started(&time.Time{}), Elapsed(time.Minute)),
		},
		{
			name:  "started null",
			json:  `{"name":"bar","duration":60000,"started":"","elapsed":1000}`,
			timer: New(Name("bar"), Duration(time.Minute), Elapsed(time.Second)),
		},
		{
			name:         "bad json",
			json:         `{"name":"bar","duration":60000,"started":"","elapsed":1000`,
			timer:        New(),
			unmarshalErr: &json.SyntaxError{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			timer := New()
			err := json.Unmarshal([]byte(test.json), &timer)
			if test.unmarshalErr == nil && err != nil {
				t.Errorf("err\nexpected: nil\n  actual: %v", err)
			} else if test.unmarshalErr != nil && err == nil {
				t.Errorf("err\nexpected: %v\n  actual: nil", test.unmarshalErr)
			} else if test.unmarshalErr == nil && err == nil {
				// no big deal
			} else {
				testErrType := reflect.TypeOf(test.unmarshalErr)
				actualErrType := reflect.TypeOf(err)
				if testErrType != actualErrType {
					t.Errorf("err type\nexpected: %v\n  actual: %v", testErrType, actualErrType)
				}
			}

			if (test.timer.started == nil) && (timer.started != nil) {
				t.Errorf("timer.started\nexpected: nil\n  actual: &%v", *timer.started)
			} else if (test.timer.started != nil) && (timer.started == nil) {
				t.Errorf("timer.started\nexpected: &%v\n  actual: nil", *test.timer.started)
			} else if (test.timer.started == nil) && (timer.started == nil) {
				// no big deal
			} else if *test.timer.started != *timer.started {
				t.Errorf("timer.started\nexpected: &%v\n  actual: &%v", *test.timer.started, *timer.started)
			}

			timer.started = nil
			test.timer.started = nil
			if timer != test.timer {
				t.Errorf("timer (ignore started)\nexpected: %v\n  actual: %v", test.timer, timer)
			}
		})
	}
}
