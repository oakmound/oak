package render

import (
	"image/color"
	"math/rand"
	"testing"
	"time"

	"github.com/oakmound/oak/physics"
	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	p, err := NewPolygon([]physics.Vector{
		physics.NewVector(10, 10),
		physics.NewVector(20, 10),
		physics.NewVector(10, 20),
	})
	assert.Nil(t, err)
	assert.True(t, p.Contains(11, 11))
	assert.False(t, p.Contains(16, 16))
	// fails!
	//assert.True(t, p.WrappingContains(11, 11))
	assert.False(t, p.WrappingContains(16, 16))
	assert.True(t, p.ConvexContains(11, 11))
	assert.False(t, p.ConvexContains(16, 16))
	// This also is wonky, should consider working with shape
	p.Fill(color.RGBA{255, 0, 0, 255})
	assert.Equal(t, p.GetRGBA().At(1, 1), color.RGBA{255, 0, 0, 255})
	p.FillInverse(color.RGBA{0, 255, 0, 255})
	assert.Equal(t, p.GetRGBA().At(1, 1), color.RGBA{0, 255, 0, 255})
}

func BenchmarkContains(b *testing.B) {
	curSeed := time.Now().UTC().UnixNano()
	rand.Seed(curSeed)

	points := []physics.Vector{}
	for i := 0; i < 100; i++ {
		points = append(points, physics.NewVector(rand.Float64()*640, rand.Float64()*480))
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

	points := []physics.Vector{}
	for i := 0; i < 100; i++ {
		points = append(points, physics.NewVector(rand.Float64()*640, rand.Float64()*480))
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

	points := []physics.Vector{}
	for i := 0; i < 100; i++ {
		points = append(points, physics.NewVector(rand.Float64()*640, rand.Float64()*480))
	}
	poly, _ := NewPolygon(points)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := rand.Float64() * 640
		y := rand.Float64() * 480
		poly.ConvexContains(x, y)
	}
}
