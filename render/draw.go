package render

import (
	"image/color"

	"time"

	"github.com/oakmound/oak/timing"
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
func LoadSpriteAndDraw(filename string, l int) (Renderable, error) {
	s, err := LoadSprite(filename)
	if err != nil {
		return nil, err
	}
	return Draw(s, l)
}

// DrawColor is equivalent to LoadSpriteAndDraw,
// but with colorboxes.
func DrawColor(c color.Color, x1, y1, x2, y2 float64, layer, stackLayer int) (Renderable, error) {
	cb := NewColorBox(int(x2), int(y2), c)
	cb.ShiftX(x1)
	cb.ShiftY(y1)
	if len(GlobalDrawStack.as) == 1 {
		return Draw(cb, layer)
	}
	return Draw(cb, stackLayer, layer)
}

// DrawForTime draws and after d undraws an element
func DrawForTime(r Renderable, l int, d time.Duration) error {
	_, err := Draw(r, l)
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
