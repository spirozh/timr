package timer

import (
	"testing"
	"time"
)

func TestOptions(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		timer := New(Name("foo"))
		if timer.Name() != "foo" {
			t.Fatal("Name set failed")
		}

		oldConfig := timer.Config(Name("bar"))
		if timer.Name() != "bar" {
			t.Fatal("Name option failed")
		}

		timer.Config(oldConfig...)
		if timer.Name() != "foo" {
			t.Fatal("Name reset failed")
		}
	})

	t.Run("duration", func(t *testing.T) {
		timer := New(Duration(0))
		if timer.duration != 0 {
			t.Fatal("Duration set failed")
		}

		oldConfig := timer.Config(Duration(time.Hour))
		if timer.duration != time.Hour {
			t.Fatal("Duration option failed")
		}

		timer.Config(oldConfig...)
		if timer.duration != 0 {
			t.Fatal("Duration reset failed")
		}
	})

	t.Run("started", func(t *testing.T) {
		timer := New(Started(nil))
		if timer.started != nil {
			t.Fatal("Started set failed")
		}

		oldConfig := timer.Config(Started(&time.Time{}))
		if *timer.started != (time.Time{}) {
			t.Fatal("Started option failed")
		}

		timer.Config(oldConfig...)
		if timer.started != nil {
			t.Fatal("Started reset failed")
		}
	})

	t.Run("elapsed", func(t *testing.T) {
		timer := New(Elapsed(0))
		if timer.elapsed != 0 {
			t.Fatal("Duration set failed")
		}

		oldConfig := timer.Config(Elapsed(time.Hour))
		if timer.elapsed != time.Hour {
			t.Fatal("Elapsed option failed")
		}

		timer.Config(oldConfig...)
		if timer.elapsed != 0 {
			t.Fatal("Elapsed reset failed")
		}
	})

	t.Run("config", func(t *testing.T) {
		now := time.Time{}
		timer := New(Name("foo"), Duration(time.Hour), Started(&now), Elapsed(time.Minute))
		b0, _ := timer.MarshalJSON()

		oldConfig := timer.Config(Name("bar"), Duration(time.Minute), Started(nil), Elapsed(time.Second))
		b1, _ := timer.MarshalJSON()
		if string(b0) == string(b1) {
			t.Fatal("setting config failed")
		}

		timer.Config(oldConfig...)
		b2, _ := timer.MarshalJSON()
		if string(b0) != string(b2) {
			t.Fatal("restoring config failed")
		}

		// original assertions
	})
}
