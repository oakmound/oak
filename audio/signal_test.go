package audio

import "testing"

func TestPosSignal(t *testing.T) {
	psgn := NewPosSignal(1, 2, 3)
	ok, x, y := psgn.GetPos()
	if !ok {
		t.Fatalf("expected getPos to return true")
	}
	if x != 2 {
		t.Fatalf("expected x of %v, got %v", 2, x)
	}
	if y != 3 {
		t.Fatalf("expected y of %v, got %v", 3, x)
	}
}
