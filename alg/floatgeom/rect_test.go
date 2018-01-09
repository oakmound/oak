package floatgeom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRectConstructors(t *testing.T) {
	assert.Equal(t, Rect2{Min: Point2{0, 0}, Max: Point2{1, 1}}, NewRect2(0, 0, 1, 1))
	assert.Equal(t, Rect2{Min: Point2{0, 0}, Max: Point2{1, 1}}, NewRect2(1, 1, 0, 0))
	assert.Equal(t, Rect2{Min: Point2{0, 0}, Max: Point2{1, 1}}, NewRect2WH(0, 0, 1, 1))
	assert.Equal(t, Rect2{Min: Point2{0, 0}, Max: Point2{1, 1}}, NewRect2WH(1, 1, -1, -1))
	assert.Equal(t, Rect2{Min: Point2{0, 0}, Max: Point2{1, 1}},
		NewBoundingRect2(Point2{0, 0}, Point2{0, 1}, Point2{1, 0}, Point2{1, 1}))
	assert.Equal(t, Rect3{Min: Point3{0, 0, 0}, Max: Point3{1, 1, 1}}, NewRect3(0, 0, 0, 1, 1, 1))
	assert.Equal(t, Rect3{Min: Point3{0, 0, 0}, Max: Point3{1, 1, 1}}, NewRect3(1, 1, 1, 0, 0, 0))
	assert.Equal(t, Rect3{Min: Point3{0, 0, 0}, Max: Point3{1, 1, 1}}, NewRect3WH(0, 0, 0, 1, 1, 1))
	assert.Equal(t, Rect3{Min: Point3{0, 0, 0}, Max: Point3{1, 1, 1}}, NewRect3WH(1, 1, 1, -1, -1, -1))
	assert.Equal(t, Rect3{Min: Point3{0, 0, 0}, Max: Point3{1, 1, 1}},
		NewBoundingRect3(Point3{0, 0, 0}, Point3{0, .5, 1}, Point3{.5, 1, 0}, Point3{1, 0, .5}))
}

func TestRectAccess(t *testing.T) {
	r2 := NewRect2(0, 1, 2, 3)
	r3 := NewRect3(0, 1, 2, 3, 4, 5)
	assert.Equal(t, 4.0, r2.Area())
	assert.Equal(t, 27.0, r3.Space())
	assert.Equal(t, 2.0, r2.W())
	assert.Equal(t, 2.0, r2.H())
	assert.Equal(t, 3.0, r3.W())
	assert.Equal(t, 3.0, r3.H())
	assert.Equal(t, 3.0, r3.D())
	assert.Equal(t, 1.0, r2.Midpoint(0))
	assert.Equal(t, 2.0, r2.Midpoint(1))
	assert.Equal(t, 1.5, r3.Midpoint(0))
	assert.Equal(t, 2.5, r3.Midpoint(1))
	assert.Equal(t, 3.5, r3.Midpoint(2))
	assert.Equal(t, 8.0, r2.Perimeter())
	assert.Equal(t, 36.0, r3.Margin())
}

func TestRect2Contains(t *testing.T) {
	r2 := NewRect2(0, 0, 10, 10)
	expected := []bool{true}

	for i, e := range expected {
		c := pt2cases[i]
		a := Point2{c.x1, c.y1}
		assert.Equal(t, e, r2.Contains(a))
	}
}

func TestRect3Contains(t *testing.T) {
	r3 := NewRect3(0, 0, 0, 10, 10, 10)
	expected := []bool{true}

	for i, e := range expected {
		c := pt3cases[i]
		a := Point3{c.x1, c.y1, c.z1}
		assert.Equal(t, e, r3.Contains(a))
	}
}

var (
	r2cases = []Rect2{
		NewRect2(1, 1, 2, 2),
		NewRect2(3, 3, 11, 11),
		NewRect2(11, 11, 12, 12),
	}
	r3cases = []Rect3{
		NewRect3(1, 1, 1, 2, 2, 2),
		NewRect3(3, 3, 3, 11, 11, 11),
		NewRect3(11, 11, 11, 12, 12, 12),
	}
)

func TestRect2ContainsRect(t *testing.T) {
	r2 := NewRect2(0, 0, 10, 10)
	expected := []bool{true, false, false}

	for i, e := range expected {
		c := r2cases[i]
		assert.Equal(t, e, r2.ContainsRect(c))
	}
}

func TestRect3ContainsRect(t *testing.T) {
	r3 := NewRect3(0, 0, 0, 10, 10, 10)
	expected := []bool{true, false, false}

	for i, e := range expected {
		c := r3cases[i]
		assert.Equal(t, e, r3.ContainsRect(c))
	}
}

func TestRect2Intersects(t *testing.T) {
	r2 := NewRect2(0, 0, 10, 10)
	expected := []bool{true, true, false}

	for i, e := range expected {
		c := r2cases[i]
		assert.Equal(t, e, r2.Intersects(c))
	}
}

func TestRect3Intersects(t *testing.T) {
	r3 := NewRect3(0, 0, 0, 10, 10, 10)
	expected := []bool{true, true, false}

	for i, e := range expected {
		c := r3cases[i]
		assert.Equal(t, e, r3.Intersects(c))
	}
}

func TestMaxRectDimensions(t *testing.T) {
	assert.Equal(t, 2, Rect2{}.MaxDimensions())
	assert.Equal(t, 3, Rect3{}.MaxDimensions())
}

func TestRect2GreaterOf(t *testing.T) {
	r2 := NewRect2(0, 0, 10, 10)
	expected := []Rect2{
		NewRect2(0, 0, 10, 10),
		NewRect2(0, 0, 11, 11),
		NewRect2(0, 0, 12, 12),
	}

	for i, e := range expected {
		c := r2cases[i]
		assert.Equal(t, e, r2.GreaterOf(c))
	}
}

func TestRect3GreaterOf(t *testing.T) {
	r3 := NewRect3(0, 0, 0, 10, 10, 10)
	expected := []Rect3{
		NewRect3(0, 0, 0, 10, 10, 10),
		NewRect3(0, 0, 0, 11, 11, 11),
		NewRect3(0, 0, 0, 12, 12, 12),
	}

	for i, e := range expected {
		c := r3cases[i]
		assert.Equal(t, e, r3.GreaterOf(c))
	}
}
