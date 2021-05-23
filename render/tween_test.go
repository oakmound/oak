package render

import (
	"image/color"
	"testing"
)

func TestTween(t *testing.T) {
	start := NewColorBox(10, 10, color.RGBA{0, 0, 0, 0})
	end := NewColorBox(10, 10, color.RGBA{255, 255, 255, 255})
	// I didn't expect to have to give frames - 2 here
	tween := Tween(start.GetRGBA(), end.GetRGBA(), 254)
	for i, rgba := range tween {
		c := rgba.At(0, 0)
		r, g, b, a := c.RGBA()
		// I mean, I can guess that this should be near 255 but
		// I had to just jump around to actually find 257 (and I've
		// had to do this before, and remember this same experience)
		v := uint32(257 * i)

		if r != v || g != v || b != v || a != v {
			t.Fatalf("tween color mismatch")
		}
	}
}
