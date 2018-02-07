package render

import (
	"image/color"
	"testing"

	"github.com/oakmound/oak/shape"
	"github.com/stretchr/testify/assert"
)

func TestSimpleBezierLine(t *testing.T) {
	bz, err := shape.BezierCurve(0, 0, 10, 10)
	assert.Nil(t, err)
	sp := BezierLine(bz, color.RGBA{255, 255, 255, 255})
	rgba := sp.GetRGBA()
	for i := 0; i < 10; i++ {
		assert.Equal(t, color.RGBA{255, 255, 255, 255}, rgba.At(i, i))
	}

	bz, err = shape.BezierCurve(10, 10, 0, 0)
	assert.Nil(t, err)
	sp = BezierLine(bz, color.RGBA{255, 255, 255, 255})
	rgba = sp.GetRGBA()
	for i := 0; i < 10; i++ {
		assert.Equal(t, color.RGBA{255, 255, 255, 255}, rgba.At(i, i))
	}
}
