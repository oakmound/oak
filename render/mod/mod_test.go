package mod

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPallete(t *testing.T) {
	w, h := 10, 10
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	ConformToPalleteFilter([]color.Color{color.RGBA{128, 0, 0, 128}})(rgba)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			assert.Equal(t, color.RGBA{128, 0, 0, 128}, rgba.At(x, y))
		}
	}

	InPlace(ConformToPallete(
		[]color.Color{color.RGBA{64, 0, 0, 128}}))(rgba)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			assert.Equal(t, color.RGBA{64, 0, 0, 128}, rgba.At(x, y))
		}
	}
}
