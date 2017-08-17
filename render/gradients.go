package render

import (
	"image/color"

	"github.com/200sc/go-dist/colorrange"
)

// GradientColorAt returns a new color via a gradient between two colors and the progress between them
func GradientColorAt(c1, c2 color.Color, progress float64) color.Color {
	return colorrange.NewLinear(c1, c2).Percentile(progress)
}
