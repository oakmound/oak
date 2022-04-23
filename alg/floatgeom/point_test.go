package floatgeom

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/oakmound/oak/v4/alg"
)

func Seed() {
	rand.Seed(time.Now().UnixNano())
}

func TestPointRotate(t *testing.T) {
	p := Point2{0, 1}
	if -1.0 != p.Rotate(90).X() {
		t.Fatalf("expected %v got %v", 0, p.Rotate(90).X())
	}
	if -1.0 != p.RotateRadians(math.Pi).Y() {
		t.Fatalf("expected %v got %v", 0, p.RotateRadians(math.Pi).Y())
	}
}

func TestPointNormalize(t *testing.T) {
	p1 := Point2{100, 200}.Normalize()
	p2 := Point3{100, 200, 300}.Normalize()
	p3 := Point4{100, 200, 300, 400}.Normalize()

	if !alg.F64eqEps(p1.X(), 1/math.Sqrt(5), .0001) {
		t.Fatalf("expected %v got %v", 1/math.Sqrt(5), p1.X())
	}
	if !alg.F64eqEps(p1.Y(), 2/math.Sqrt(5), .0001) {
		t.Fatalf("expected %v got %v", 2/math.Sqrt(5), p1.Y())
	}
	if !alg.F64eqEps(p2.X(), 1/math.Sqrt(14), .0001) {
		t.Fatalf("expected %v got %v", 1/math.Sqrt(14), p2.X())
	}
	if !alg.F64eqEps(p2.Y(), 2/math.Sqrt(14), .0001) {
		t.Fatalf("expected %v got %v", 2/math.Sqrt(14), p2.Y())
	}
	if !alg.F64eqEps(p2.Z(), 3/math.Sqrt(14), .0001) {
		t.Fatalf("expected %v got %v", 3/math.Sqrt(14), p2.Z())
	}
	if !alg.F64eqEps(p3.W(), 1/math.Sqrt(30), .0001) {
		t.Fatalf("expected %v got %v", 1/math.Sqrt(30), p3.W())
	}
	if !alg.F64eqEps(p3.X(), 2/math.Sqrt(30), .0001) {
		t.Fatalf("expected %v got %v", 2/math.Sqrt(30), p3.X())
	}
	if !alg.F64eqEps(p3.Y(), 3/math.Sqrt(30), .0001) {
		t.Fatalf("expected %v got %v", 3/math.Sqrt(30), p3.Y())
	}
	if !alg.F64eqEps(p3.Z(), 4/math.Sqrt(30), .0001) {
		t.Fatalf("expected %v got %v", 4/math.Sqrt(30), p3.Z())
	}

	p4 := Point2{0, 0}
	p5 := Point3{0, 0, 0}
	p6 := Point4{0, 0, 0, 0}
	if p4 != p4.Normalize() {
		t.Fatalf("expected %v got %v", 4, p4.Normalize())
	}
	if p5 != p5.Normalize() {
		t.Fatalf("expected %v got %v", 5, p5.Normalize())
	}
	if p6 != p6.Normalize() {
		t.Fatalf("expected %v got %v", 6, p6.Normalize())
	}
}

func TestPointProject(t *testing.T) {
	Seed()
	for i := 0; i < randTests; i++ {
		x, y, z := rand.Float64(), rand.Float64(), rand.Float64()
		p := Point3{x, y, z}
		if p.ProjectZ() != (Point2{x, y}) {
			t.Fatalf("expected %v got %v", p.ProjectZ(), Point2{x, y})
		}
		if p.ProjectY() != (Point2{x, z}) {
			t.Fatalf("expected %v got %v", p.ProjectY(), Point2{x, z})
		}
		if p.ProjectX() != (Point2{y, z}) {
			t.Fatalf("expected %v got %v", p.ProjectX(), Point2{y, z})
		}
	}
}

func TestCrossProduct(t *testing.T) {
	p1 := Point3{1, 2, 1}
	p2 := Point3{3, 1, 3}

	if p1.Cross(p2) != (Point3{5, 0, -5}) {
		t.Fatalf("expected %v got %v", p1.Cross(p2), Point3{5, 0, -5})
	}
}

func TestPointConstMods(t *testing.T) {
	p1 := Point2{1, 1}
	p2 := Point3{1, 1, 1}
	p3 := Point4{1, 1, 1, 1}
	if p1.MulConst(5) != (Point2{5, 5}) {
		t.Fatalf("expected %v got %v", p1.MulConst(5), Point2{5, 5})
	}
	if p2.MulConst(100) != (Point3{100, 100, 100}) {
		t.Fatalf("expected %v got %v", p2.MulConst(100), Point3{100, 100, 100})
	}
	if p3.MulConst(500) != (Point4{500, 500, 500, 500}) {
		t.Fatalf("expected %v got %v", p3.MulConst(500), Point4{500, 500, 500, 500})
	}
	p4 := Point2{2, 2}
	p5 := Point3{2, 2, 2}
	p6 := Point4{2, 2, 2, 2}
	if p4.DivConst(4) != (Point2{.5, .5}) {
		t.Fatalf("expected %v got %v", p4.DivConst(4), Point2{.5, .5})
	}
	if p5.DivConst(8) != (Point3{.25, .25, .25}) {
		t.Fatalf("expected %v got %v", p5.DivConst(8), Point3{.25, .25, .25})
	}
	if p6.DivConst(2) != (Point4{1, 1, 1, 1}) {
		t.Fatalf("expected %v got %v", p6.DivConst(2), Point4{1, 1, 1, 1})
	}
}

func TestAnglePoints(t *testing.T) {
	p1 := Point2{1, 0}
	p2 := Point2{0, 1}

	if p1 != AnglePoint(0) {
		t.Fatalf("expected %v got %v", 1, AnglePoint(0))
	}
	if p1 != RadianPoint(0) {
		t.Fatalf("expected %v got %v", 1, RadianPoint(0))
	}
	p3 := AnglePoint(90)
	p4 := RadianPoint(math.Pi / 2)
	if !alg.F64eqEps(p2.Y(), p3.Y(), .0001) {
		t.Fatalf("expected %v got %v", p3.Y(), p2.Y())
	}
	if !alg.F64eqEps(p2.Y(), p4.Y(), .0001) {
		t.Fatalf("expected %v got %v", p4.Y(), p2.Y())
	}

	deg := p1.ToAngle()
	rds := p1.ToRadians()

	if 0 != deg {
		t.Fatalf("expected %v got %v", 0, deg)
	}
	if 0 != rds {
		t.Fatalf("expected %v got %v", 0, rds)
	}

	deg = p2.ToAngle()
	rds = p2.ToRadians()

	if !alg.F64eqEps(90.0, deg, .0001) {
		t.Fatalf("expected %v got %v", 90.0, deg)
	}
	if !alg.F64eqEps(math.Pi/2, rds, .0001) {
		t.Fatalf("expected %v got %v", math.Pi/2, rds)
	}

	p5 := Point2{-1, 0}

	if !alg.F64eqEps(45.0, p2.AngleTo(p5), .0001) {
		t.Fatalf("expected %v got %v", 45.0, p2.AngleTo(p5))
	}
	if !alg.F64eqEps(math.Pi/4, p2.RadiansTo(p5), .0001) {
		t.Fatalf("expected %v got %v", math.Pi/4, p2.RadiansTo(p5))
	}
}

func TestPointGreaterOf(t *testing.T) {
	a := Point2{0, 1}
	b := Point2{1, 0}
	if a.GreaterOf(b) != (Point2{1, 1}) {
		t.Fatalf("expected %v got %v", a.GreaterOf(b), Point2{1, 1})
	}

	c := Point3{0, 1, 2}
	d := Point3{1, 2, 0}
	if c.GreaterOf(d) != (Point3{1, 2, 2}) {
		t.Fatalf("expected %v got %v", c.GreaterOf(d), Point3{1, 2, 2})
	}
}

func TestPointLesserOf(t *testing.T) {
	a := Point2{0, 1}
	b := Point2{1, 0}
	if a.LesserOf(b) != (Point2{0, 0}) {
		t.Fatalf("expected %v got %v", a.LesserOf(b), Point2{0, 0})
	}

	c := Point3{0, 1, 2}
	d := Point3{1, 2, 0}
	if c.LesserOf(d) != (Point3{0, 1, 0}) {
		t.Fatalf("expected %v got %v", c.LesserOf(d), Point3{0, 1, 0})
	}
}

func TestPointAccess(t *testing.T) {
	a := Point4{0, 1, 2, 3}
	if 0 != a.Dim(0) {
		t.Fatalf("expected %v got %v", 0, a.Dim(0))
	}
	if 1 != a.Dim(1) {
		t.Fatalf("expected %v got %v", 0, a.Dim(1))
	}
	if 2 != a.Dim(2) {
		t.Fatalf("expected %v got %v", 0, a.Dim(2))
	}
	if 3 != a.Dim(3) {
		t.Fatalf("expected %v got %v", 0, a.Dim(3))
	}
	if 0 != a.W() {
		t.Fatalf("expected %v got %v", 0, a.W())
	}
	if 1 != a.X() {
		t.Fatalf("expected %v got %v", 0, a.X())
	}
	if 2 != a.Y() {
		t.Fatalf("expected %v got %v", 0, a.Y())
	}
	if 3 != a.Z() {
		t.Fatalf("expected %v got %v", 0, a.Z())
	}

	b := Point3{0, 1, 2}
	if 0 != b.Dim(0) {
		t.Fatalf("expected %v got %v", 0, b.Dim(0))
	}
	if 1 != b.Dim(1) {
		t.Fatalf("expected %v got %v", 0, b.Dim(1))
	}
	if 2 != b.Dim(2) {
		t.Fatalf("expected %v got %v", 0, b.Dim(2))
	}
	if 0 != b.X() {
		t.Fatalf("expected %v got %v", 0, b.X())
	}
	if 1 != b.Y() {
		t.Fatalf("expected %v got %v", 0, b.Y())
	}
	if 2 != b.Z() {
		t.Fatalf("expected %v got %v", 0, b.Z())
	}

	c := Point2{0, 1}
	if 0 != c.Dim(0) {
		t.Fatalf("expected %v got %v", 0, c.Dim(0))
	}
	if 1 != c.Dim(1) {
		t.Fatalf("expected %v got %v", 0, c.Dim(1))
	}
	if 0 != c.X() {
		t.Fatalf("expected %v got %v", 0, c.X())
	}
	if 1 != c.Y() {
		t.Fatalf("expected %v got %v", 0, c.Y())
	}
}

// Pattern here: there's a set of input pairs here
// each test takes these and has expected outputs for each pair index.
var (
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
		if a.Distance(b) != e {
			t.Fatalf("expected %v got %v", a.Distance(b), e)
		}
	}
}

func TestPointDistance2(t *testing.T) {
	expected := []float64{math.Sqrt(2)}

	for i, e := range expected {
		c := pt2cases[i]
		a := Point2{c.x1, c.y1}
		b := Point2{c.x2, c.y2}
		if a.Distance(b) != e {
			t.Fatalf("expected %v got %v", a.Distance(b), e)
		}
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
		if a.Add(b) != e {
			t.Fatalf("expected %v got %v", a.Add(b), e)
		}
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
		if a.Add(b) != e {
			t.Fatalf("expected %v got %v", a.Add(b), e)
		}
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
		if a.Sub(b) != e {
			t.Fatalf("expected %v got %v", a.Sub(b), e)
		}
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
		if a.Sub(b) != e {
			t.Fatalf("expected %v got %v", a.Sub(b), e)
		}
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
		if a.Mul(b) != e {
			t.Fatalf("expected %v got %v", a.Mul(b), e)
		}
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
		if a.Mul(b) != e {
			t.Fatalf("expected %v got %v", a.Mul(b), e)
		}
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
		if a.Div(b) != e {
			t.Fatalf("expected %v got %v", a.Div(b), e)
		}
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
		if a.Div(b) != e {
			t.Fatalf("expected %v got %v", a.Div(b), e)
		}
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

		if expected3 != (Point3{x, y, z}).ToRect(span) {
			t.Fatalf("expected %v got %v", expected3, (Point3{x, y, z}).ToRect(span))
		}

		if expected2 != (Point2{x, y}).ToRect(span) {
			t.Fatalf("expected %v got %v", expected2, (Point2{x, y}).ToRect(span))
		}
	}
}

func TestQuaternionMultiplication(t *testing.T) {
	a := Point4{1, 1, 1, 1}.Normalize()
	b := a.Inverse()
	c := Point4{1, 0, 0, 0}
	if c != a.MulQuat(b) {
		t.Fatalf("expected %v got %v", c, a.MulQuat(b))
	}
}
