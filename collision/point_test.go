package collision

import (
	"testing"
)

func TestNewPoint(t *testing.T) {
	p := NewPoint(nil, 10, 10)
	if p.X() != 10 {
		t.Fatalf("bad x, expected %v got %v", 10, p.X())
	}
	if p.Y() != 10 {
		t.Fatalf("bad y, expected %v got %v", 10, p.Y())
	}
	if !p.IsNil() {
		t.Fatalf("nil point should have been nil")
	}
	p2 := NewPoint(&Space{}, 0, 0)
	if p2.IsNil() {
		t.Fatalf("set point should not have been nil")
	}
}
