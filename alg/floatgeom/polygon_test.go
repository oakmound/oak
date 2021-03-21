package floatgeom

import (
	"math/rand"
	"testing"
	"time"
)

func TestPolygon_Contains(t *testing.T) {
	p := NewPolygon2(
		Point2{10, 10},
		Point2{20, 10},
		Point2{10, 20},
	)
	if !p.Contains(11, 11) {
		t.Fatalf("polygon did not contain 11,11")
	}
	if p.Contains(16, 16) {
		t.Fatalf("polygon contained 16,16")
	}
	if p.Contains(40, 40) {
		t.Fatalf("polygon contained 40,40")
	}
	if !p.ConvexContains(11, 11) {
		t.Fatalf("convex polygon did not contain 11,11")
	}
	if p.ConvexContains(16, 16) {
		t.Fatalf("convex polygon contained 16,16")
	}
	if p.ConvexContains(40, 40) {
		t.Fatalf("convex polygon contained 40,40")
	}
}

func TestPolygon_getSide(t *testing.T) {
	v := Point2{2, 2}
	if getSide(v, v) != 0 {
		t.Fatalf("getSide did not confirm identical points")
	}
	p := NewPolygon2(
		Point2{0, 0},
		Point2{0, 10},
		Point2{10, 10},
		Point2{10, 0},
	)
	if p.ConvexContains(0, 5) {
		t.Fatalf("point in line with polygon was contained")
	}
}

var benchContains bool

func BenchmarkPolygonContains(b *testing.B) {
	curSeed := time.Now().UTC().UnixNano()
	rand.Seed(curSeed)

	points := []Point2{}
	for i := 0; i < 100; i++ {
		points = append(points, Point2{rand.Float64() * 640, rand.Float64() * 480})
	}
	poly := NewPolygon2(points[0], points[1], points[2], points[3:]...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := rand.Float64() * 640
		y := rand.Float64() * 480
		benchContains = poly.Contains(x, y)
	}
}

func BenchmarkPolygonConvexContains(b *testing.B) {
	curSeed := time.Now().UTC().UnixNano()
	rand.Seed(curSeed)

	points := []Point2{}
	for i := 0; i < 100; i++ {
		points = append(points, Point2{rand.Float64() * 640, rand.Float64() * 480})
	}
	poly := NewPolygon2(points[0], points[1], points[2], points[3:]...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := rand.Float64() * 640
		y := rand.Float64() * 480
		benchContains = poly.ConvexContains(x, y)
	}
}
