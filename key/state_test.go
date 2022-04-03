package key

import (
	"testing"
	"time"
)

func TestState(t *testing.T) {
	ks := NewState()
	ks.SetDown(A)
	if !ks.IsDown(A) {
		t.Fatalf("a was not set down")
	}
	ks.SetUp(A)
	if ks.IsDown(A) {
		t.Fatalf("a was not set up")
	}
	ks.SetDown(A)
	time.Sleep(2 * time.Second)
	ok, d := ks.IsHeld(A)
	if !ok {
		t.Fatalf("a was not held down")
	}
	if d < 2000*time.Millisecond {
		t.Fatalf("a was not held down for sleep length")
	}
	ks.SetUp(A)
	ok, d = ks.IsHeld(A)
	if ok {
		t.Fatalf("a was not released")
	}
	if d != 0 {
		t.Fatalf("a hold was not reset")
	}
}
