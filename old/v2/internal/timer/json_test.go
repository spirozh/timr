package timer_test

import (
	"encoding/json"
	"reflect"
	"spirozh/timr/internal/timer"
	"testing"
	"time"
)

func Test_JSONMarshal(t *testing.T) {
	type test struct {
		name  string
		timer any
		json  string
	}

	tests := []test{
		{
			name:  "zero timer",
			timer: timer.New(),
			json:  `{"name":"","duration":0,"started":null,"elapsed":0}`,
		},
		{
			name:  "zero started",
			timer: timer.New(timer.Name("foo"), timer.Duration(time.Second), timer.Started(&time.Time{}), timer.Elapsed(time.Minute)),
			json:  `{"name":"foo","duration":1000,"started":"0001-01-01T00:00:00.0000","elapsed":60000}`,
		},
		{
			name:  "nil started",
			timer: timer.New(timer.Name("bar"), timer.Duration(time.Minute), timer.Started(nil), timer.Elapsed(time.Second)),
			json:  `{"name":"bar","duration":60000,"started":null,"elapsed":1000}`,
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
		unmarshalErr error
		timer        any
	}

	tests := []test{
		{
			name:  "started zero",
			json:  `{"name":"foo","duration":1000,"started":"0001-01-01T00:00:00.0000","elapsed":60000}`,
			timer: timer.New(timer.Name("foo"), timer.Duration(time.Second), timer.Started(&time.Time{}), timer.Elapsed(time.Minute)),
		},
		{
			name:  "started null",
			json:  `{"name":"bar","duration":60000,"started":null,"elapsed":1000}`,
			timer: timer.New(timer.Name("bar"), timer.Duration(time.Minute), timer.Elapsed(time.Second)),
		},
		{
			name:         "bad json (bad date)",
			json:         `{"name":"bar","duration":60000,"started":"not a date","elapsed":1000}`,
			timer:        timer.New(),
			unmarshalErr: &time.ParseError{},
		},
		{
			name:         "bad json (missing closing brace)",
			json:         `{"name":"bar","duration":60000,"started":null,"elapsed":1000`,
			timer:        timer.New(),
			unmarshalErr: &json.SyntaxError{},
		},
		{
			name:         "bad json (numeric name)",
			json:         `{"name":0,"duration":60000,"started":null,"elapsed":1000}`,
			timer:        timer.New(),
			unmarshalErr: &json.UnmarshalTypeError{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			timer := timer.New()
			err := json.Unmarshal([]byte(test.json), &timer)
			if expected, actual := reflect.TypeOf(test.unmarshalErr), reflect.TypeOf(err); expected != actual {
				t.Errorf("err type while unmarshalling:\nexpected: %v\n  actual: %v", expected, actual)
			}

			if !reflect.DeepEqual(test.timer, timer) {
				t.Errorf("timer\nexpected: %s\n  actual: %s", test.timer, timer)
			}
		})
	}
}
