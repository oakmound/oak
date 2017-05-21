package render

import (
	"image"
	"image/color"
)

// Tween takes two images and returns a set of images tweening
// between the two over some number of frames
func Tween(a image.Image, b image.Image, frames int) []*image.RGBA {
	bounds := a.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y

	tweened := make([]*image.RGBA, frames+2)
	progress := 0.0
	inc := 1.0 / float64(len(tweened)-1)
	for i := range tweened {
		tweened[i] = image.NewRGBA(image.Rect(0, 0, w, h))
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				r1, g1, b1, a1 := a.At(x, y).RGBA()
				r2, g2, b2, a2 := b.At(x, y).RGBA()

				r1f := float64(r1)
				g1f := float64(g1)
				b1f := float64(b1)
				a1f := float64(a1)

				r2f := float64(r2)
				g2f := float64(g2)
				b2f := float64(b2)
				a2f := float64(a2)

				r1f *= 1 - progress
				g1f *= 1 - progress
				b1f *= 1 - progress
				a1f *= 1 - progress

				r2f *= progress
				g2f *= progress
				b2f *= progress
				a2f *= progress

				c := color.RGBA64{uint16(r1f + r2f), uint16(g1f + g2f),
					uint16(b1f + b2f), uint16(a1f + a2f)}

				tweened[i].Set(x, y, c)
			}
		}
		progress += inc
	}

	return tweened
}
