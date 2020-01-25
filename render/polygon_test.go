package render

import (
	"image/color"
	"math/rand"
	"testing"
	"time"

	"github.com/oakmound/oak/v2/alg/floatgeom"
	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	p, err := NewPolygon(
		floatgeom.Point2{10, 10},
		floatgeom.Point2{20, 10},
		floatgeom.Point2{10, 20},
	)
	assert.Nil(t, err)
	assert.True(t, p.Contains(11, 11))
	assert.False(t, p.Contains(16, 16))
	assert.False(t, p.Contains(40, 40))

	assert.True(t, p.ConvexContains(11, 11))
	assert.False(t, p.ConvexContains(16, 16))
	// This also is wonky, should consider working with shape
	p.Fill(color.RGBA{255, 0, 0, 255})
	assert.Equal(t, p.GetRGBA().At(1, 1), color.RGBA{255, 0, 0, 255})
	p.FillInverse(color.RGBA{0, 255, 0, 255})
	assert.Equal(t, p.GetRGBA().At(1, 1), color.RGBA{0, 255, 0, 255})
}

func TestCrossEqual(t *testing.T) {
	v := floatgeom.Point2{2, 2}
	assert.Equal(t, 0, getSide(v, v))
	p, err := NewPolygon(
		floatgeom.Point2{0, 0},
		floatgeom.Point2{0, 10},
		floatgeom.Point2{10, 10},
		floatgeom.Point2{10, 0},
	)
	assert.Nil(t, err)
	assert.False(t, p.ConvexContains(0, 5))
}

func TestPolygonFns(t *testing.T) {
	p, err := NewPolygon(
		floatgeom.Point2{0, 0},
		floatgeom.Point2{0, 10},
		floatgeom.Point2{10, 10},
		floatgeom.Point2{10, 0},
	)
	assert.Nil(t, err)
	cmp := p.GetOutline(color.RGBA{255, 0, 0, 255})
	assert.NotNil(t, cmp)
	assert.Len(t, cmp.rs, 4)

	err = p.UpdatePoints(
		floatgeom.Point2{0, 0},
	)
	assert.NotNil(t, err)

	err = p.UpdatePoints(
		floatgeom.Point2{0, 0},
		floatgeom.Point2{10, 0},
		floatgeom.Point2{0, 10},
	)
	assert.Nil(t, err)

	_, err = NewPolygon()
	assert.NotNil(t, err)
}

func TestStrictPolygon(t *testing.T) {
	p, err := NewStrictPolygon(floatgeom.NewRect2(0, 0, 640, 480),
		floatgeom.Point2{0, 0},
		floatgeom.Point2{10, 0},
		floatgeom.Point2{0, 10})
	assert.Nil(t, err)
	assert.Equal(t, 0.0, p.X())
	assert.Equal(t, 0.0, p.Y())
	w, h := p.GetDims()
	assert.Equal(t, 640, w)
	assert.Equal(t, 480, h)

	_, err = NewStrictPolygon(floatgeom.NewRect2(0, 0, 640, 480))
	assert.NotNil(t, err)
}

func BenchmarkContains(b *testing.B) {
	curSeed := time.Now().UTC().UnixNano()
	rand.Seed(curSeed)

	points := []floatgeom.Point2{}
	for i := 0; i < 100; i++ {
		points = append(points, floatgeom.Point2{rand.Float64() * 640, rand.Float64() * 480})
	}
	poly, _ := NewPolygon(points...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := rand.Float64() * 640
		y := rand.Float64() * 480
		poly.Contains(x, y)
	}
}

func BenchmarkConvexContains(b *testing.B) {
	curSeed := time.Now().UTC().UnixNano()
	rand.Seed(curSeed)

	points := []floatgeom.Point2{}
	for i := 0; i < 100; i++ {
		points = append(points, floatgeom.Point2{rand.Float64() * 640, rand.Float64() * 480})
	}
	poly, _ := NewPolygon(points...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := rand.Float64() * 640
		y := rand.Float64() * 480
		poly.ConvexContains(x, y)
	}
}
