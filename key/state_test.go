package key

import (
	"testing"
	"time"
)

func TestState(t *testing.T) {
	ks := NewState()
	ks.SetDown("Test")
	if !ks.IsDown("Test") {
		t.Fatalf("test was not set down")
	}
	ks.SetUp("Test")
	if ks.IsDown("Test") {
		t.Fatalf("test was not set up")
	}
	ks.SetDown("Test")
	time.Sleep(2 * time.Second)
	ok, d := ks.IsHeld("Test")
	if !ok {
		t.Fatalf("test was not held down")
	}
	if d < 2000*time.Millisecond {
		t.Fatalf("test was not held down for sleep length")
	}
	ks.SetUp("Test")
	ok, d = ks.IsHeld("Test")
	if ok {
		t.Fatalf("test was not released")
	}
	if d != 0 {
		t.Fatalf("test hold was not reset")
	}
}
