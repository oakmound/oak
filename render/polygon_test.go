package render

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkContains(b *testing.B) {
	curSeed := time.Now().UTC().UnixNano()
	rand.Seed(curSeed)

	points := []Point{}
	for i := 0; i < 100; i++ {
		points = append(points, Point{rand.Float64() * 640, rand.Float64() * 480})
	}
	poly, _ := NewPolygon(points)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := rand.Float64() * 640
		y := rand.Float64() * 480
		poly.Contains(x, y)
	}
}

func BenchmarkWrappingContains(b *testing.B) {
	curSeed := time.Now().UTC().UnixNano()
	rand.Seed(curSeed)

	points := []Point{}
	for i := 0; i < 100; i++ {
		points = append(points, Point{rand.Float64() * 640, rand.Float64() * 480})
	}
	poly, _ := NewPolygon(points)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := rand.Float64() * 640
		y := rand.Float64() * 480
		poly.WrappingContains(x, y)
	}
}

func BenchmarkConvexContains(b *testing.B) {
	curSeed := time.Now().UTC().UnixNano()
	rand.Seed(curSeed)

	points := []Point{}
	for i := 0; i < 100; i++ {
		points = append(points, Point{rand.Float64() * 640, rand.Float64() * 480})
	}
	poly, _ := NewPolygon(points)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := rand.Float64() * 640
		y := rand.Float64() * 480
		poly.ConvexContains(x, y)
	}
}
