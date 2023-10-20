package timer_test

import (
	"encoding/json"
	"spirozh/timr/internal/timer"
	"testing"
	"time"
)

func TestOptions(t *testing.T) {

	t.Run("name", func(t *testing.T) {
		timr := timer.New(timer.Name("foo"))
		if timr.Name() != "foo" {
			t.Fatal(
				"Name set failed",
				"\nexpected: \"foo\"",
				"\n  actual: ", timr.Name(),
			)
		}

		restoreConfig := timr.Config(timer.Name("bar"))
		if timr.Name() != "bar" {
			t.Fatal(
				"Name option failed",
				"\nexpected: \"bar\"",
				"\n  actual: ", timr.Name(),
			)
		}

		timr.Config(restoreConfig...)
		if timr.Name() != "foo" {
			t.Fatal(
				"Name restore failed",
				"\nexpected: \"foo\"",
				"\n  actual: ", timr.Name(),
			)
		}
	})

	t.Run("duration", func(t *testing.T) {
		timr := timer.New(timer.Duration(time.Hour))
		if actual := timr.Remaining(time.Time{}); actual != time.Hour {
			t.Fatal(
				"Duration set failed",
				"\nexpected: ", 0,
				"\n  actual: ", actual,
			)
		}

		restoreConfig := timr.Config(timer.Duration(0))
		if actual := timr.Remaining(time.Time{}); actual != 0 {
			t.Fatal(
				"Duration option failed",
				"\nexpected: ", time.Hour,
				"\n  actual: ", actual,
			)
		}

		timr.Config(restoreConfig...)
		if actual := timr.Remaining(time.Time{}); actual != time.Hour {
			t.Fatal(
				"Duration restore failed",
				"\nexpected: ", 0,
				"\n  actual: ", actual,
			)
		}
	})

	t.Run("started", func(t *testing.T) {
		timr := timer.New(timer.Started(&time.Time{}))
		if timr.Elapsed(time.Time{}.Add(time.Minute)) != time.Minute {
			t.Fatal("Started set failed")
		}

		restoreConfig := timr.Config(timer.Started(nil))
		if timr.Elapsed(time.Time{}.Add(time.Minute)) != 0 {
			t.Fatal("Started option failed")
		}

		timr.Config(restoreConfig...)
		if timr.Elapsed(time.Time{}.Add(time.Minute)) != time.Minute {
			t.Fatal("Started restore failed")
		}
	})

	t.Run("elapsed", func(t *testing.T) {
		timr := timer.New(timer.Elapsed(time.Hour))
		if actual := timr.Remaining(time.Time{}); actual != -time.Hour {
			t.Fatal("Elapsed set failed: ", actual)
		}

		restoreConfig := timr.Config(timer.Elapsed(0))
		if timr.Remaining(time.Time{}) != 0 {
			t.Fatal("Elapsed option failed")
		}

		timr.Config(restoreConfig...)
		if actual := timr.Remaining(time.Time{}); actual != -time.Hour {
			t.Fatal("Elapsed restore failed: ", actual)
		}
	})

	t.Run("config", func(t *testing.T) {
		timr := timer.New(timer.Name("foo"), timer.Duration(time.Hour), timer.Started(&time.Time{}), timer.Elapsed(time.Minute))

		b0, _ := json.Marshal(timr)

		oldConfig := timr.Config(timer.Name("bar"), timer.Duration(time.Minute), timer.Started(nil), timer.Elapsed(time.Second))
		b1, _ := json.Marshal(timr)
		if string(b0) == string(b1) {
			t.Fatal("setting config failed")
		}

		timr.Config(oldConfig...)
		b2, _ := json.Marshal(timr)
		if string(b0) != string(b2) {
			t.Fatal("restoring config failed")
		}
	})
}
