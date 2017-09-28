package render

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleCircle(t *testing.T) {
	rgba := image.NewRGBA(image.Rect(0, 0, 10, 10))
	DrawCircle(rgba, color.RGBA{255, 255, 255, 255}, 5, 0, 0, 0)
	// For better or for worse, the current implementation produces
	// . . . . . . . . . .
	// . . . x x x x x . .
	// . . x x       x x .
	// . x x           x x
	// . x               x
	// . x               x
	// . x               x
	// . x x           x x
	// . . x x       x x .
	// . . . x x x x x . .
	// This should change in the future, probably leaning towards using Beziers.
	boolExpected := [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 1, 1, 1, 1, 0, 0},
		{0, 0, 1, 1, 0, 0, 0, 1, 1, 0},
		{0, 1, 1, 0, 0, 0, 0, 0, 1, 1},
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
		{0, 1, 1, 0, 0, 0, 0, 0, 1, 1},
		{0, 0, 1, 1, 0, 0, 0, 1, 1, 0},
		{0, 0, 0, 1, 1, 1, 1, 1, 0, 0},
	}
	for x, col := range boolExpected {
		for y, b := range col {
			if b == 0 {
				assert.Equal(t, color.RGBA{0, 0, 0, 0}, rgba.At(x, y))
			} else {
				assert.Equal(t, color.RGBA{255, 255, 255, 255}, rgba.At(x, y))
			}
		}
	}
}
