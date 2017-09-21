package shape

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	randTestCt = 100
)

func TestBezier(t *testing.T) {
	bp := BezierPoint{2, 5}
	for i := 0; i < randTestCt; i++ {
		x, y := bp.Pos(rand.Float64())
		assert.Equal(t, 2.0, x)
		assert.Equal(t, 5.0, y)
	}

	bn := BezierNode{
		BezierPoint{0, 0},
		BezierPoint{0, 1},
	}

	for i := 0; i < randTestCt; i++ {
		r := rand.Float64()
		x, y := bn.Pos(r)
		assert.Equal(t, 0.0, x)
		assert.Equal(t, r, y)
	}

	bc, err := BezierCurve(0, 0, 0, 1)

	assert.Equal(t, bn, bc)
	assert.Nil(t, err)

	_, err = BezierCurve()
	assert.NotNil(t, err)

	_, err = BezierCurve(0, 0, 0)
	assert.NotNil(t, err)
}
