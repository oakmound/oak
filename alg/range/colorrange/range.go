// Package colorrange provides distributions that accept and return color.Colors.
package colorrange

import (
	"image/color"
)

// Range represents a range of colors
type Range interface {
	Poll() color.Color
	EnforceRange(color.Color) color.Color
	Percentile(f float64) color.Color
}
