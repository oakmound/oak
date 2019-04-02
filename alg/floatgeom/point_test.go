package floatgeom

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Seed() {
	rand.Seed(time.Now().UnixNano())
}

func TestPointRotate(t *testing.T) {
	p := Point2{0, 1}
	assert.Equal(t, -1.0, p.Rotate(90).X())
	assert.Equal(t, -1.0, p.RotateRadians(math.Pi).Y())
}

func TestPointNormalize(t *testing.T) {
	p1 := Point2{100, 200}.Normalize()
	p2 := Point3{100, 200, 300}.Normalize()
	p3 := Point4{100, 200, 300, 400}.Normalize()

	assert.InEpsilon(t, p1.X(), 1/math.Sqrt(5), .0001)
	assert.InEpsilon(t, p1.Y(), 2/math.Sqrt(5), .0001)
	assert.InEpsilon(t, p2.X(), 1/math.Sqrt(14), .0001)
	assert.InEpsilon(t, p2.Y(), 2/math.Sqrt(14), .0001)
	assert.InEpsilon(t, p2.Z(), 3/math.Sqrt(14), .0001)
	assert.InEpsilon(t, p3.W(), 1/math.Sqrt(30), .0001)
	assert.InEpsilon(t, p3.X(), 2/math.Sqrt(30), .0001)
	assert.InEpsilon(t, p3.Y(), 3/math.Sqrt(30), .0001)
	assert.InEpsilon(t, p3.Z(), 4/math.Sqrt(30), .0001)

	p4 := Point2{0, 0}
	p5 := Point3{0, 0, 0}
	p6 := Point4{0, 0, 0, 0}
	assert.Equal(t, p4, p4.Normalize())
	assert.Equal(t, p5, p5.Normalize())
	assert.Equal(t, p6, p6.Normalize())
}

func TestPointProject(t *testing.T) {
	Seed()
	for i := 0; i < randTests; i++ {
		x, y, z := rand.Float64(), rand.Float64(), rand.Float64()
		p := Point3{x, y, z}
		assert.Equal(t, Point2{x, y}, p.ProjectZ())
		assert.Equal(t, Point2{x, z}, p.ProjectY())
		assert.Equal(t, Point2{y, z}, p.ProjectX())
	}
}

func TestCrossProduct(t *testing.T) {
	p1 := Point3{1, 2, 1}
	p2 := Point3{3, 1, 3}

	assert.Equal(t, Point3{5, 0, -5}, p1.Cross(p2))
}

func TestPointConstMods(t *testing.T) {
	p1 := Point2{1, 1}
	p2 := Point3{1, 1, 1}
	p3 := Point4{1, 1, 1, 1}
	assert.Equal(t, Point2{5, 5}, p1.MulConst(5))
	assert.Equal(t, Point3{100, 100, 100}, p2.MulConst(100))
	assert.Equal(t, Point4{500, 500, 500, 500}, p3.MulConst(500))
	p4 := Point2{2, 2}
	p5 := Point3{2, 2, 2}
	p6 := Point4{2, 2, 2, 2}
	assert.Equal(t, Point2{.5, .5}, p4.DivConst(4))
	assert.Equal(t, Point3{.25, .25, .25}, p5.DivConst(8))
	assert.Equal(t, Point4{1, 1, 1, 1}, p6.DivConst(2))
}

func TestAnglePoints(t *testing.T) {
	p1 := Point2{1, 0}
	p2 := Point2{0, 1}

	assert.Equal(t, p1, AnglePoint(0))
	assert.Equal(t, p1, RadianPoint(0))
	p3 := AnglePoint(90)
	p4 := RadianPoint(math.Pi / 2)
	assert.InEpsilon(t, p2.Y(), p3.Y(), .0001)
	assert.InEpsilon(t, p2.Y(), p4.Y(), .0001)

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

func TestPointAccess(t *testing.T) {
	a := Point4{0, 1, 2, 3}
	assert.Equal(t, 0.0, a.Dim(0))
	assert.Equal(t, 1.0, a.Dim(1))
	assert.Equal(t, 2.0, a.Dim(2))
	assert.Equal(t, 3.0, a.Dim(3))
	assert.Equal(t, 0.0, a.W())
	assert.Equal(t, 1.0, a.X())
	assert.Equal(t, 2.0, a.Y())
	assert.Equal(t, 3.0, a.Z())

	b := Point3{0, 1, 2}
	assert.Equal(t, 0.0, b.Dim(0))
	assert.Equal(t, 1.0, b.Dim(1))
	assert.Equal(t, 2.0, b.Dim(2))
	assert.Equal(t, 0.0, b.X())
	assert.Equal(t, 1.0, b.Y())
	assert.Equal(t, 2.0, b.Z())

	c := Point2{0, 1}
	assert.Equal(t, 0.0, c.Dim(0))
	assert.Equal(t, 1.0, c.Dim(1))
	assert.Equal(t, 0.0, c.X())
	assert.Equal(t, 1.0, c.Y())
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

func TestQuaternionMultiplication(t *testing.T) {
	a := Point4{1, 1, 1, 1}.Normalize()
	b := a.Inverse()
	c := Point4{1, 0, 0, 0}
	assert.Equal(t, c, a.MulQuat(b))
}
