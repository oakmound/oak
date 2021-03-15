package floatgeom

import (
	"testing"
)

func TestTriangleNormal(t *testing.T) {
	a := Tri3{
		Point3{0, 0, 0},
		Point3{1, 0, 0},
		Point3{0, 1, 0},
	}
	e := Point3{0, 0, 1}
	if e != a.Normal() {
		t.Fatalf("expected %v got %v", e, a.Normal())
	}
}

func TestTriangleBarycentric(t *testing.T) {
	a := Tri3{
		Point3{0, 0, 0},
		Point3{1, 0, 0},
		Point3{0, 1, 0},
	}
	e := Point3{1, 1, -1}
	if e != a.Barycentric(1, 1) {
		t.Fatalf("expected %v got %v", e, a.Barycentric(1, 1))
	}
	e = Point3{0.5, 0.5, 0}
	if e != a.Barycentric(.5, .5) {
		t.Fatalf("expected %v got %v", e, a.Barycentric(.5, .5))
	}
}
