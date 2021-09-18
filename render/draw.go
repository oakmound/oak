package render

import (
	"image/color"
)

var (
	// emptyRenderable is a simple renderable that can be used
	// for pseudo-nil renderables that need to be something
	emptyRenderable = NewColorBox(1, 1, color.RGBA{0, 0, 0, 0})
)

// EmptyRenderable returns a minimal, 1-width and height pseudo-nil
// Renderable (and Modifiable)
func EmptyRenderable() Modifiable {
	return emptyRenderable.Copy()
}

// DrawColor is equivalent to LoadSpriteAndDraw,
// but with colorboxes.
func DrawColor(c color.Color, x, y, w, h float64, layers ...int) (Renderable, error) {
	cb := NewColorBox(int(w), int(h), c)
	cb.ShiftX(x)
	cb.ShiftY(y)
	return Draw(cb, layers...)
}

// DrawPoint draws a color on the screen as a single-widthed
// pixel (box)
func DrawPoint(c color.Color, x1, y1 float64, layers ...int) (Renderable, error) {
	return DrawColor(c, x1, y1, 1, 1, layers...)
}
