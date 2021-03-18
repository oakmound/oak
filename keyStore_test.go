package oak

import (
	"testing"
	"time"
)

func TestKeyStore(t *testing.T) {
	// Todo: maybe not strings, it's debatable what is more useful
	// Strings
	// Pros: Easy to write
	// Cons: Can write wrong thing, "Space" instead of "Spacebar"
	//
	// Enum
	// Pros: Uses less space, probably less time as well
	// Cons: Requires import, key.A instead of "A", keybinds require an extended const block
	SetDown("Test")
	if !IsDown("Test") {
		t.Fatalf("test was not set down")
	}
	SetUp("Test")
	if IsDown("Test") {
		t.Fatalf("test was not set up")
	}
	SetDown("Test")
	time.Sleep(2 * time.Second)
	ok, d := IsHeld("Test")
	if !ok {
		t.Fatalf("test was not held down")
	}
	if d < 2000*time.Millisecond {
		t.Fatalf("test was not held down for sleep length")
	}
	SetUp("Test")
	ok, d = IsHeld("Test")
	if ok {
		t.Fatalf("test was not released")
	}
	if d != 0 {
		t.Fatalf("test hold was not reset")
	}

	// KeyBind
	if GetKeyBind("Test") != "Test" {
		t.Fatalf("getKeyBind did not return identiy for non-bound key")
	}
	BindKey("Test", "Bound")
	if GetKeyBind("Test") != "Bound" {
		t.Fatalf("getKeyBind did not return bound value for bound key")
	}
}
