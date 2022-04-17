package oak

import (
	"image"
	"image/color"
	"testing"
)

func TestScreenFilter(t *testing.T) {
	c1 := NewWindow()
	blackAndWhite := color.Palette{
		color.RGBA{0, 0, 0, 255},
		color.RGBA{255, 255, 255, 255},
	}
	c1.SetPalette(blackAndWhite)
	buf := image.NewRGBA(image.Rect(0, 0, 1, 1))
	c1.prePublish(buf)
}
