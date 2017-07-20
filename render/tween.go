package render

import (
	"image"
	"image/color"
)

// Tween takes two images and returns a set of images tweening
// between the two over some number of frames
func Tween(start image.Image, end image.Image, frames int) []*image.RGBA {
	bounds := start.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y

	tweened := make([]*image.RGBA, frames+2)
	progress := 0.0
	inc := 1.0 / float64(len(tweened)-1)
	for i := range tweened {
		tweened[i] = image.NewRGBA(image.Rect(0, 0, w, h))
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				r1, g1, b1, a1 := start.At(x, y).RGBA()
				r2, g2, b2, a2 := end.At(x, y).RGBA()

				r1f := float64(r1) * (1 - progress)
				g1f := float64(g1) * (1 - progress)
				b1f := float64(b1) * (1 - progress)
				a1f := float64(a1) * (1 - progress)

				r2f := float64(r2) * progress
				g2f := float64(g2) * progress
				b2f := float64(b2) * progress
				a2f := float64(a2) * progress

				c := color.RGBA64{uint16(r1f + r2f), uint16(g1f + g2f),
					uint16(b1f + b2f), uint16(a1f + a2f)}

				tweened[i].Set(x, y, c)
			}
		}
		progress += inc
	}

	return tweened
}
