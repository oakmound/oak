package render

import (
	"testing"

	"github.com/akavel/polyclip-go"
	"github.com/stretchr/testify/assert"
)

func TestDrawPolygon(t *testing.T) {
	rh := RenderableHeap{}

	r := rh.DrawPolygonDim()
	assert.Equal(t, 0.0, r.Min.X())
	assert.Equal(t, 0.0, r.Min.Y())
	assert.Equal(t, 0.0, r.Max.X())
	assert.Equal(t, 0.0, r.Max.Y())

	x := 10.0
	y := 10.0
	x2 := 20.0
	y2 := 20.0

	pgn := polyclip.Polygon{{{X: x, Y: y}, {X: x, Y: y2}, {X: x2, Y: y2}, {X: x2, Y: y}}}
	rh.SetDrawPolygon(pgn)

	r = rh.DrawPolygonDim()
	assert.Equal(t, x, r.Min.X())
	assert.Equal(t, y, r.Min.Y())
	assert.Equal(t, x2, r.Max.X())
	assert.Equal(t, y2, r.Max.Y())

	type testcase struct {
		elems         [4]int
		shouldSucceed bool
	}

	tests := []testcase{
		{[4]int{0, 0, 0, 0}, false},
		{[4]int{0, 0, 30, 30}, true},
		{[4]int{15, 15, 17, 17}, true},
	}

	for _, cas := range tests {
		assert.Equal(t, cas.shouldSucceed, rh.InDrawPolygon(cas.elems[0], cas.elems[1], cas.elems[2], cas.elems[3]))
	}

	rh.ClearDrawPolygon()

	for _, cas := range tests {
		assert.Equal(t, true, rh.InDrawPolygon(cas.elems[0], cas.elems[1], cas.elems[2], cas.elems[3]))

	}
}
