package render

import (
	"image/color"
)

//GradientColorAt returns a new color via a gradient between two colors and the progress between them
func GradientColorAt(c1, c2 color.Color, progress float64) color.RGBA64 {
	r, g, b, a := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return color.RGBA64{
		uint16OnScale(r, r2, progress),
		uint16OnScale(g, g2, progress),
		uint16OnScale(b, b2, progress),
		uint16OnScale(a, a2, progress),
	}
}
