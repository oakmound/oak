package collision

import (
	"testing"

	"github.com/oakmound/oak/v2/alg/floatgeom"

	"github.com/oakmound/oak/v2/physics"
	"github.com/stretchr/testify/assert"
)

func TestSpaceFuncs(t *testing.T) {
	Clear()
	s := NewUnassignedSpace(10, 10, 10, 10)
	assert.NotEmpty(t, s.String())
	if s.W() != 10.0 {
		t.Fatalf("expected 10 width, got %v", s.W())
	}
	if s.H() != 10.0 {
		t.Fatalf("expected 10 height, got %v", s.H())
	}

	// Getters
	cx, cy := s.GetCenter()
	assert.Equal(t, cx, float64(15))
	assert.Equal(t, cy, float64(15))
	x, y := s.GetPos()
	assert.Equal(t, x, float64(10))
	assert.Equal(t, y, float64(10))

	// Positional comparison
	s2 := NewUnassignedSpace(20, 20, 10, 10)
	assert.True(t, s2.Above(s) < 0)
	assert.True(t, s2.Below(s) > 0)
	assert.True(t, s2.LeftOf(s) < 0)
	assert.True(t, s2.RightOf(s) > 0)

	// Containment
	assert.False(t, s2.Contains(s))
	s3 := NewUnassignedSpace(5, 5, 20, 20)
	assert.True(t, s3.Contains(s))
	s4 := NewUnassignedSpace(15, 15, 10, 10)

	// Overlap
	xover, yover := s4.Overlap(s)
	assert.Equal(t, xover, -5.0)
	assert.Equal(t, yover, -5.0)
	xover, yover = s.Overlap(s4)
	assert.Equal(t, xover, 5.0)
	assert.Equal(t, yover, 5.0)
	xover, yover = s.Overlap(s2)
	assert.Equal(t, xover, 0.0)
	assert.Equal(t, yover, 0.0)
	ov := s.OverlapVector(s4)
	assert.Equal(t, ov, physics.NewVector(5, 5))
	spaces := s.SubtractRect(1, 1, 8, 8)
	assert.Equal(t, len(spaces), 4)
}

func TestNewRect(t *testing.T) {
	s := NewUnassignedSpace(0, 0, 0, 0)
	assert.Equal(t, 1.0, s.GetW())
	assert.Equal(t, 1.0, s.GetH())
	s = NewUnassignedSpace(0, 0, -10, -10)
	assert.Equal(t, 10.0, s.GetW())
	assert.Equal(t, 10.0, s.GetH())
	s = NewRectSpace(floatgeom.NewRect3WH(0, 0, 0, 10, 10, 0), 0, 0)
	assert.Equal(t, 10.0, s.GetW())
	assert.Equal(t, 10.0, s.GetH())
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
