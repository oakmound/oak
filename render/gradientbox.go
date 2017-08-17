package render

import (
	"image"
	"image/color"

	"github.com/200sc/go-dist/colorrange"
)

// NewGradientBox returns a gradient box defined on the two input colors
// and the given progress function
func NewGradientBox(w, h int, startColor, endColor color.Color, pFunction Progress) *Sprite {
	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)

	crange := colorrange.NewLinear(startColor, endColor)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			progress := pFunction(x, y, w, h)
			c := crange.Percentile(progress)
			rgba.Set(x, y, c)
		}
	}
	return NewSprite(0, 0, rgba)
}

// NewHorizontalGradientBox returns a gradient box with a horizontal gradient from
// the start to end color, left to right.
func NewHorizontalGradientBox(w, h int, startColor, endColor color.Color) *Sprite {
	return NewGradientBox(w, h, startColor, endColor, HorizontalProgress)
}

// NewVerticalGradientBox returns a gradient box with a vertical gradient from
// the start to end color, top to bottom.
func NewVerticalGradientBox(w, h int, startColor, endColor color.Color) *Sprite {
	return NewGradientBox(w, h, startColor, endColor, VerticalProgress)
}

// NewCircularGradientBox returns a gradient box where the center will be startColor
// and the gradient will radiate as a circle out from the center.
func NewCircularGradientBox(w, h int, startColor, endColor color.Color) *Sprite {
	return NewGradientBox(w, h, startColor, endColor, CircularProgress)
}
