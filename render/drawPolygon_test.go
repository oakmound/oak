package render

import (
	"testing"

	"github.com/akavel/polyclip-go"
	"github.com/stretchr/testify/assert"
)

func TestDrawPolygon(t *testing.T) {

	a, b, c, d := DrawPolygonDim()
	assert.Equal(t, 0, a)
	assert.Equal(t, 0, b)
	assert.Equal(t, 0, c)
	assert.Equal(t, 0, d)

	x := 10.0
	y := 10.0
	x2 := 20.0
	y2 := 20.0

	pgn := polyclip.Polygon{{{X: x, Y: y}, {X: x, Y: y2}, {X: x2, Y: y2}, {X: x2, Y: y}}}
	SetDrawPolygon(pgn)

	x3, y3, x4, y4 := DrawPolygonDim()
	assert.Equal(t, int(x), x3)
	assert.Equal(t, int(y), y3)
	assert.Equal(t, int(x2), x4)
	assert.Equal(t, int(y2), y4)

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
		assert.Equal(t, cas.shouldSucceed, InDrawPolygon(cas.elems[0], cas.elems[1], cas.elems[2], cas.elems[3]))
	}

	ClearDrawPolygon()

	for _, cas := range tests {
		assert.Equal(t, true, InDrawPolygon(cas.elems[0], cas.elems[1], cas.elems[2], cas.elems[3]))

	}
}
