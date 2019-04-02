package intgeom

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Seed() {
	rand.Seed(time.Now().UnixNano())
}

func TestPointProject(t *testing.T) {
	Seed()
	for i := 0; i < randTests; i++ {
		x, y, z := rand.Intn(100), rand.Intn(100), rand.Intn(100)
		p := Point3{x, y, z}
		assert.Equal(t, Point2{x, y}, p.ProjectZ())
		assert.Equal(t, Point2{x, z}, p.ProjectY())
		assert.Equal(t, Point2{y, z}, p.ProjectX())
	}
}

func TestPointConstMods(t *testing.T) {
	p1 := Point2{1, 1}
	p2 := Point3{1, 1, 1}
	assert.Equal(t, Point2{5, 5}, p1.MulConst(5))
	assert.Equal(t, Point3{100, 100, 100}, p2.MulConst(100))
	p3 := Point2{2, 2}
	p4 := Point3{2, 2, 2}
	assert.Equal(t, Point2{0, 0}, p3.DivConst(4))
	assert.Equal(t, Point3{0, 0, 0}, p4.DivConst(8))
}

func TestAnglePoints(t *testing.T) {
	p1 := Point2{1, 0}
	p2 := Point2{0, 1}

	deg := p1.ToAngle()
	rds := p1.ToRadians()

	assert.Equal(t, 0.0, deg)
	assert.Equal(t, 0.0, rds)

	deg = p2.ToAngle()
	rds = p2.ToRadians()

	assert.InEpsilon(t, 90.0, deg, .0001)
	assert.InEpsilon(t, math.Pi/2, rds, .0001)

	p5 := Point2{-1, 0}

	assert.InEpsilon(t, 45.0, p2.AngleTo(p5), .0001)
	assert.InEpsilon(t, math.Pi/4, p2.RadiansTo(p5), .0001)
}

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

func TestPointMagnitude(t *testing.T) {
	a := Point2{2, 2}
	b := Point2{1, 1}
	c := a.LesserOf(b)
	assert.True(t, a.Magnitude() > c.Magnitude())
	assert.True(t, b.Magnitude() == c.Magnitude())

	d := Point3{2, 3, 4}
	e := Point3{1, 2, 3}
	f := d.LesserOf(e)
	fmt.Println(d.Magnitude(), e.Magnitude(), f.Magnitude())

	assert.True(t, d.Magnitude() > f.Magnitude())
	assert.True(t, e.Magnitude() == f.Magnitude())
}

func TestPointAccess(t *testing.T) {
	a := Point3{0, 1, 2}
	assert.Equal(t, 0, a.Dim(0))
	assert.Equal(t, 1, a.Dim(1))
	assert.Equal(t, 2, a.Dim(2))
	assert.Equal(t, 0, a.X())
	assert.Equal(t, 1, a.Y())
	assert.Equal(t, 2, a.Z())

	b := Point2{0, 1}
	assert.Equal(t, 0, b.Dim(0))
	assert.Equal(t, 1, b.Dim(1))
	assert.Equal(t, 0, b.X())
	assert.Equal(t, 1, b.Y())
}

// Pattern here: there's a set of input pairs here
// each test takes these and has expected outputs for each pair index.
var (
	// Todo: add more test cases
	pt3cases = []struct{ x1, y1, z1, x2, y2, z2 int }{
		{0, 0, 0, 1, 1, 1},
	}
	pt2cases = []struct{ x1, y1, x2, y2 int }{
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
		span := rand.Intn(100)
		x := rand.Intn(100)
		y := rand.Intn(100)
		z := rand.Intn(100)

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
