package test

import (
	"fmt"
	"testing"
)

func NotEqual(t *testing.T, notExpected, actual any, msgAndArgs ...any) {
	t.Helper()
	if actual == notExpected {
		msg := "Failure\n"
		message := messageFrom(msgAndArgs...)
		if len(message) > 0 {
			msg += fmt.Sprintln(message)
		}
		msg += fmt.Sprintf("\nexpected value not to be: %v", notExpected)
		t.Error(msg)
	}
}

func Equal(t *testing.T, expected, actual any, msgAndArgs ...any) {
	t.Helper()
	if expected != actual {
		msg := "Failure\n"
		message := messageFrom(msgAndArgs...)
		if len(message) > 0 {
			msg += fmt.Sprintln(message)
		}
		msg += fmt.Sprintf("\nexpected: %v", expected)
		msg += fmt.Sprintf("\nactual  : %v", actual)
		t.Error(msg)
	}
}

func ElementsMatch(t *testing.T, args ...any) {
	t.Helper()
	msg := "Args:"
	for i, arg := range args {
		msg += fmt.Sprintf("\n\targ(%d): '%v'", i, arg)
	}

	if len(args) == 0 {
		msg = "No args"
	}
	t.Error("TODO: ElementsMatch not implemented.\n" + msg)

}

func messageFrom(msgAndArgs ...any) string {
	return "TODO: messageFrom(msgAndArgs...) not implemented."
}
