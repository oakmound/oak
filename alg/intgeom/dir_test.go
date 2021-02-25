package intgeom

import "testing"

func TestDirMethods(t *testing.T) {
	d := Dir2{10, 12}
	if d.X() != 10 {
		t.Fatalf("expected 10 for x, got %v", d.X())
	}
	if d.Y() != 12 {
		t.Fatalf("expected 12 for y, got %v", d.Y())
	}
}
