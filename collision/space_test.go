package collision

import (
	"testing"

	"github.com/oakmound/oak/v3/alg/floatgeom"

	"github.com/oakmound/oak/v3/physics"
)

func TestSpaceFuncs(t *testing.T) {
	Clear()
	s := NewUnassignedSpace(10, 10, 10, 10)
	if s.W() != 10.0 {
		t.Fatalf("expected 10 width, got %v", s.W())
	}
	if s.H() != 10.0 {
		t.Fatalf("expected 10 height, got %v", s.H())
	}

	// Getters
	cx, cy := s.GetCenter()
	if cx != float64(15) {
		t.Fatalf("expected %v got %v", cx, float64(15))
	}
	if cy != float64(15) {
		t.Fatalf("expected %v got %v", cy, float64(15))
	}
	x, y := s.GetPos()
	if x != float64(10) {
		t.Fatalf("expected %v got %v", x, float64(10))
	}
	if y != float64(10) {
		t.Fatalf("expected %v got %v", y, float64(10))
	}

	// Positional comparison
	s2 := NewUnassignedSpace(20, 20, 10, 10)
	if !(s2.Above(s) < 0) {
		t.Fatalf("s2 should not be above s")
	}
	if !(s2.Below(s) > 0) {
		t.Fatalf("s2 should be below s")
	}
	if !(s2.LeftOf(s) < 0) {
		t.Fatalf("s2 should not be left of s")
	}
	if !(s2.RightOf(s) > 0) {
		t.Fatalf("s2 should be right of s")
	}

	// Containment
	if s2.Contains(s) {
		t.Fatalf("s2 should not contain s")
	}
	s3 := NewUnassignedSpace(5, 5, 20, 20)
	if !s3.Contains(s) {
		t.Fatalf("s3 should contain s")
	}
	s4 := NewUnassignedSpace(15, 15, 10, 10)

	// Overlap
	xover, yover := s4.Overlap(s)
	if xover != -5.0 {
		t.Fatalf("expected %v got %v", xover, -5.0)
	}
	if yover != -5.0 {
		t.Fatalf("expected %v got %v", yover, -5.0)
	}
	xover, yover = s.Overlap(s4)
	if xover != 5.0 {
		t.Fatalf("expected %v got %v", xover, 5.0)
	}
	if yover != 5.0 {
		t.Fatalf("expected %v got %v", yover, 5.0)
	}
	xover, yover = s.Overlap(s2)
	if xover != 0.0 {
		t.Fatalf("expected %v got %v", xover, 0.0)
	}
	if yover != 0.0 {
		t.Fatalf("expected %v got %v", yover, 0.0)
	}
	ov := s.OverlapVector(s4)
	if ov.X() != 5 || ov.Y() != 5 {
		t.Fatalf("expected %v got %v", ov, physics.NewVector(5, 5))
	}
	spaces := s.SubtractRect(1, 1, 8, 8)
	if len(spaces) != 4 {
		t.Fatalf("expected %v got %v", len(spaces), 4)
	}
}

func TestNewRect(t *testing.T) {
	s := NewUnassignedSpace(0, 0, 0, 0)
	if 1.0 != s.GetW() {
		t.Fatalf("expected %v got %v", 1.0, s.GetW())
	}
	if 1.0 != s.GetH() {
		t.Fatalf("expected %v got %v", 1.0, s.GetH())
	}
	s = NewUnassignedSpace(0, 0, -10, -10)
	if 10.0 != s.GetW() {
		t.Fatalf("expected %v got %v", 10.0, s.GetW())
	}
	if 10.0 != s.GetH() {
		t.Fatalf("expected %v got %v", 10.0, s.GetH())
	}
	s = NewRectSpace(floatgeom.NewRect3WH(0, 0, 0, 10, 10, 0), 0, 0)
	if 10.0 != s.GetW() {
		t.Fatalf("expected %v got %v", 10.0, s.GetW())
	}
	if 10.0 != s.GetH() {
		t.Fatalf("expected %v got %v", 10.0, s.GetH())
	}
}

func TestNewRect2Space(t *testing.T) {
	s1 := NewRect2Space(floatgeom.Rect2{
		Min: floatgeom.Point2{5, 10},
		Max: floatgeom.Point2{10, 15},
	}, 0)
	s2 := NewUnassignedSpace(5, 10, 5, 5)
	if s1.X() != s2.X() {
		t.Fatalf("mismatched X: %v vs %v", s1.X(), s2.X())
	}
	if s1.Y() != s2.Y() {
		t.Fatalf("mismatched Y: %v vs %v", s1.Y(), s2.Y())
	}
	if s1.W() != s2.W() {
		t.Fatalf("mismatched W: %v vs %v", s1.W(), s2.W())
	}
	if s1.H() != s2.H() {
		t.Fatalf("mismatched H: %v vs %v", s1.H(), s2.H())
	}
}
