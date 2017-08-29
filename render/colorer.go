package render

import "image/color"

// A Colorer takes some notion of linear progress and returns a color
type Colorer func(float64) color.Color

// IdentityColorer returns the same color it was given at initialization,
// regardless of progress.
func IdentityColorer(c color.Color) Colorer {
	return func(float64) color.Color {
		return c
	}
}
