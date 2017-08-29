package mod

import (
	"image/color"
	"testing"

	"github.com/oakmound/oak/render"
	"github.com/stretchr/testify/assert"
)

func TestPallete(t *testing.T) {
	r := render.NewColorBox(10, 10, color.RGBA{255, 0, 0, 255})
	rgba := r.GetRGBA()
	ConformToPalleteInPlace([]color.Color{color.RGBA{128, 0, 0, 128}})(rgba)
	w, h := r.GetDims()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			assert.Equal(t, color.RGBA{128, 0, 0, 128}, rgba.At(x, y))
		}
	}

	InPlace(render.ConformToPallete(
		[]color.Color{color.RGBA{64, 0, 0, 128}}))(rgba)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			assert.Equal(t, color.RGBA{64, 0, 0, 128}, rgba.At(x, y))
		}
	}
}
