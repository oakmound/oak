package floatgeom

import "testing"

func TestDirMethods(t *testing.T) {
	d := Dir2{10.0, 12.0}
	if d.X() != 10.0 {
		t.Fatalf("expected 10 for x, got %v", d.X())
	}
	if d.Y() != 12.0 {
		t.Fatalf("expected 12 for y, got %v", d.Y())
	}
}
