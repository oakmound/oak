package intgeom

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

func TestPointProject(t *testing.T) {
	Seed()
	for i := 0; i < randTests; i++ {
		x, y, z := rand.Intn(100), rand.Intn(100), rand.Intn(100)
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

func TestPointConstMods(t *testing.T) {
	p1 := Point2{1, 1}
	p2 := Point3{1, 1, 1}
	if p1.MulConst(5) != (Point2{5, 5}) {
		t.Fatalf("expected %v got %v", p1.MulConst(5), Point2{5, 5})
	}
	if p2.MulConst(100) != (Point3{100, 100, 100}) {
		t.Fatalf("expected %v got %v", p2.MulConst(100), Point3{100, 100, 100})
	}
	p3 := Point2{2, 2}
	p4 := Point3{2, 2, 2}
	if p3.DivConst(4) != (Point2{0, 0}) {
		t.Fatalf("expected %v got %v", p3.DivConst(4), Point2{0, 0})
	}
	if p4.DivConst(8) != (Point3{0, 0, 0}) {
		t.Fatalf("expected %v got %v", p4.DivConst(8), Point3{0, 0, 0})
	}
}

func TestAnglePoints(t *testing.T) {
	p1 := Point2{1, 0}
	p2 := Point2{0, 1}

	deg := p1.ToAngle()
	rds := p1.ToRadians()

	if (0.0) != (deg) {
		t.Fatalf("expected %v got %v", 0.0, deg)
	}
	if (0.0) != (rds) {
		t.Fatalf("expected %v got %v", 0.0, rds)
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
	if (a.GreaterOf(b)) != (Point2{1, 1}) {
		t.Fatalf("expected %v got %v", a.GreaterOf(b), Point2{1, 1})
	}

	c := Point3{0, 1, 2}
	d := Point3{1, 2, 0}
	if (c.GreaterOf(d)) != (Point3{1, 2, 2}) {
		t.Fatalf("expected %v got %v", c.GreaterOf(d), Point3{1, 2, 2})
	}
}

func TestPointLesserOf(t *testing.T) {
	a := Point2{0, 1}
	b := Point2{1, 0}
	if (a.LesserOf(b)) != (Point2{0, 0}) {
		t.Fatalf("expected %v got %v", a.LesserOf(b), Point2{0, 0})
	}

	c := Point3{0, 1, 2}
	d := Point3{1, 2, 0}
	if (c.LesserOf(d)) != (Point3{0, 1, 0}) {
		t.Fatalf("expected %v got %v", c.LesserOf(d), Point3{0, 1, 0})
	}
}

func TestPointMagnitude(t *testing.T) {
	a := Point2{2, 2}
	b := Point2{1, 1}
	c := a.LesserOf(b)
	if !(a.Magnitude() > c.Magnitude()) {
		t.Fatalf("a's Magnitude did not exceed c's magnitude")
	}
	if !(b.Magnitude() == c.Magnitude()) {
		t.Fatalf("b's Magnitude did equal c's magnitude")
	}

	d := Point3{2, 3, 4}
	e := Point3{1, 2, 3}
	f := d.LesserOf(e)

	if !(d.Magnitude() > f.Magnitude()) {
		t.Fatalf("d's Magnitude did not exceed f's magnitude")

	}
	if !(e.Magnitude() == f.Magnitude()) {
		t.Fatalf("e's Magnitude did not equal f's magnitude")
	}
}

func TestPointAccess(t *testing.T) {
	a := Point3{0, 1, 2}
	if (0) != (a.Dim(0)) {
		t.Fatalf("expected %v got %v", 0, a.Dim(0))
	}
	if (1) != (a.Dim(1)) {
		t.Fatalf("expected %v got %v", 1, a.Dim(1))
	}
	if (2) != (a.Dim(2)) {
		t.Fatalf("expected %v got %v", 2, a.Dim(2))
	}
	if (0) != (a.X()) {
		t.Fatalf("expected %v got %v", 0, a.X())
	}
	if (1) != (a.Y()) {
		t.Fatalf("expected %v got %v", 1, a.Y())
	}
	if (2) != (a.Z()) {
		t.Fatalf("expected %v got %v", 2, a.Z())
	}

	b := Point2{0, 1}
	if (0) != (b.Dim(0)) {
		t.Fatalf("expected %v got %v", 0, b.Dim(0))
	}
	if (1) != (b.Dim(1)) {
		t.Fatalf("expected %v got %v", 1, b.Dim(1))
	}
	if (0) != (b.X()) {
		t.Fatalf("expected %v got %v", 0, b.X())
	}
	if (1) != (b.Y()) {
		t.Fatalf("expected %v got %v", 1, b.Y())
	}
}

// Pattern here: there's a set of input pairs here
// each test takes these and has expected outputs for each pair index.
var (
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
		if (a.Distance(b)) != (e) {
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
		if (a.Distance(b)) != (e) {
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
		if (a.Add(b)) != (e) {
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
		if (a.Add(b)) != (e) {
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
		if (a.Sub(b)) != (e) {
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
		if (a.Sub(b)) != (e) {
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
		if (a.Mul(b)) != (e) {
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
		if (a.Mul(b)) != (e) {
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
		if (a.Div(b)) != (e) {
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
		if (a.Div(b)) != (e) {
			t.Fatalf("expected %v got %v", a.Div(b), e)
		}
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

		if expected3 != (Point3{x, y, z}).ToRect(span) {
			t.Fatalf("expected %v got %v", expected3, (Point3{x, y, z}).ToRect(span))
		}

		if expected2 != (Point2{x, y}).ToRect(span) {
			t.Fatalf("expected %v got %v", expected2, (Point2{x, y}).ToRect(span))
		}
	}
}
