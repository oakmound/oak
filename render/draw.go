package render

import (
	"image/color"

	"time"

	"github.com/oakmound/oak/v2/timing"
)

var (
	resetDraw bool
	// EmptyRenderable is a simple renderable that can be used
	// for pseudo-nil renderables that need to be something
	emptyRenderable = NewColorBox(1, 1, color.RGBA{0, 0, 0, 0})
)

// EmptyRenderable returns a minimal, 1-width and height pseudo-nil
// Renderable (and Modifiable)
func EmptyRenderable() Modifiable {
	return emptyRenderable.Copy()
}

// LoadSpriteAndDraw is shorthand for LoadSprite
// followed by Draw.
func LoadSpriteAndDraw(filename string, layers ...int) (Renderable, error) {
	s, err := LoadSprite(dir, filename)
	if err != nil {
		return nil, err
	}
	return Draw(s, layers...)
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

// DrawForTime draws and after d undraws an element
func DrawForTime(r Renderable, d time.Duration, layers ...int) error {
	_, err := Draw(r, layers...)
	if err != nil {
		return err
	}
	go func(r Renderable, d time.Duration) {
		timing.DoAfter(d, func() {
			r.Undraw()
		})
	}(r, d)
	return nil
}
