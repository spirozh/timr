package test

import (
	"sort"
	"testing"
)

func NotEqual(t *testing.T, notExpected, actual any, msgAndArgs ...any) {
	t.Helper()
	if actual == notExpected {
		t.Errorf("Failure\nexpected value not to be: %v", notExpected)
	}
}

func Equal(t *testing.T, expected, actual any, msgAndArgs ...any) {
	t.Helper()
	if expected != actual {
		t.Errorf("Failure\nexpected: %v\nactual  : %v", expected, actual)
	}
}

func ElementsMatch(t *testing.T, expected, actual []string) {
	t.Helper()
	if !slicesEqual(expected, actual) {
		t.Errorf("Failure:\nexpected: %#v\nactual  : %#v", expected, actual)
	}
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	c, d := make([]string, len(a)), make([]string, len(b))
	copy(c, a)
	copy(d, b)

	sort.Strings(c)
	sort.Strings(d)

	for i := range c {
		if c[i] != d[i] {
			return false
		}
	}
	return true
}
