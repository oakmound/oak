package floatgeom

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointGreaterOf(t *testing.T) {
	a := Point2{0, 1}
	b := Point2{1, 0}
	assert.Equal(t, a.GreaterOf(b), Point2{1, 1})

	c := Point3{0, 1, 2}
	d := Point3{1, 2, 0}
	assert.Equal(t, c.GreaterOf(d), Point3{1, 2, 2})
}

func TestPointLesserOf(t *testing.T) {
	a := Point2{0, 1}
	b := Point2{1, 0}
	assert.Equal(t, a.LesserOf(b), Point2{0, 0})

	c := Point3{0, 1, 2}
	d := Point3{1, 2, 0}
	assert.Equal(t, c.LesserOf(d), Point3{0, 1, 0})
}

func TestPointAccess(t *testing.T) {
	a := Point3{0, 1, 2}
	assert.Equal(t, 0.0, a.Dim(0))
	assert.Equal(t, 1.0, a.Dim(1))
	assert.Equal(t, 2.0, a.Dim(2))
	assert.Equal(t, 0.0, a.X())
	assert.Equal(t, 1.0, a.Y())
	assert.Equal(t, 2.0, a.Z())

	b := Point2{0, 1}
	assert.Equal(t, 0.0, b.Dim(0))
	assert.Equal(t, 1.0, b.Dim(1))
	assert.Equal(t, 0.0, b.X())
	assert.Equal(t, 1.0, b.Y())
}

// Pattern here: there's a set of input pairs here
// each test takes these and has expected outputs for each pair index.
var (
	// Todo: add more test cases
	pt3cases = []struct{ x1, y1, z1, x2, y2, z2 float64 }{
		{0, 0, 0, 1, 1, 1},
	}
	pt2cases = []struct{ x1, y1, x2, y2 float64 }{
		{0, 0, 1, 1},
	}
)

func TestPointDistance3(t *testing.T) {

	expected := []float64{math.Sqrt(3)}

	for i, e := range expected {
		c := pt3cases[i]
		a := Point3{c.x1, c.y1, c.z1}
		b := Point3{c.x2, c.y2, c.z2}
		assert.Equal(t, a.Distance(b), e)
	}
}

func TestPointDistance2(t *testing.T) {
	expected := []float64{math.Sqrt(2)}

	for i, e := range expected {
		c := pt2cases[i]
		a := Point2{c.x1, c.y1}
		b := Point2{c.x2, c.y2}
		assert.Equal(t, a.Distance(b), e)
	}
}

func TestPointAdd3(t *testing.T) {
	expected := []Point3{
		{1, 1, 1},
	}
	for i, e := range expected {
		c := pt3cases[i]
		a := Point3{c.x1, c.y1, c.z1}
		b := Point3{c.x2, c.y2, c.z2}
		assert.Equal(t, a.Add(b), e)
	}
}

func TestPointAdd2(t *testing.T) {
	expected := []Point2{
		{1, 1},
	}
	for i, e := range expected {
		c := pt2cases[i]
		a := Point2{c.x1, c.y1}
		b := Point2{c.x2, c.y2}
		assert.Equal(t, a.Add(b), e)
	}
}

func TestPointSub3(t *testing.T) {
	expected := []Point3{
		{-1, -1, -1},
	}
	for i, e := range expected {
		c := pt3cases[i]
		a := Point3{c.x1, c.y1, c.z1}
		b := Point3{c.x2, c.y2, c.z2}
		assert.Equal(t, a.Sub(b), e)
	}
}

func TestPointSub2(t *testing.T) {
	expected := []Point2{
		{-1, -1},
	}
	for i, e := range expected {
		c := pt2cases[i]
		a := Point2{c.x1, c.y1}
		b := Point2{c.x2, c.y2}
		assert.Equal(t, a.Sub(b), e)
	}
}

func TestPointMul3(t *testing.T) {
	expected := []Point3{
		{0, 0, 0},
	}
	for i, e := range expected {
		c := pt3cases[i]
		a := Point3{c.x1, c.y1, c.z1}
		b := Point3{c.x2, c.y2, c.z2}
		assert.Equal(t, a.Mul(b), e)
	}
}

func TestPointMul2(t *testing.T) {
	expected := []Point2{
		{0, 0},
	}
	for i, e := range expected {
		c := pt2cases[i]
		a := Point2{c.x1, c.y1}
		b := Point2{c.x2, c.y2}
		assert.Equal(t, a.Mul(b), e)
	}
}

func TestPointDiv3(t *testing.T) {
	expected := []Point3{
		{0, 0, 0},
	}
	for i, e := range expected {
		c := pt3cases[i]
		a := Point3{c.x1, c.y1, c.z1}
		b := Point3{c.x2, c.y2, c.z2}
		assert.Equal(t, a.Div(b), e)
	}
}

func TestPointDiv2(t *testing.T) {
	expected := []Point2{
		{0, 0},
	}
	for i, e := range expected {
		c := pt2cases[i]
		a := Point2{c.x1, c.y1}
		b := Point2{c.x2, c.y2}
		assert.Equal(t, a.Div(b), e)
	}
}

var (
	randTests = 100
)

func TestPointToRec(t *testing.T) {
	for i := 0; i < randTests; i++ {
		span := rand.Float64()
		x := rand.Float64()
		y := rand.Float64()
		z := rand.Float64()

		expected3 := Rect3{Min: Point3{x, y, z}, Max: Point3{x + span, y + span, z + span}}
		expected2 := Rect2{Min: Point2{x, y}, Max: Point2{x + span, y + span}}

		assert.Equal(t,
			expected3,
			Point3{x, y, z}.ToRect(span),
		)

		assert.Equal(t,
			expected2,
			Point2{x, y}.ToRect(span),
		)
	}
}
