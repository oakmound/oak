package render

import (
	"image/color"

	"time"

	"github.com/oakmound/oak/dlog"
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
	s := LoadSprite(filename)
	return Draw(s, l)
}

// DrawColor is equivalent to LoadSpriteAndDraw,
// but with colorboxes.
func DrawColor(c color.Color, x1, y1, x2, y2 float64, layer, stackLayer int) Renderable {
	cb := NewColorBox(int(x2), int(y2), c)
	cb.ShiftX(x1)
	cb.ShiftY(y1)
	if len(GlobalDrawStack.as) == 1 {
		_, err := Draw(cb, layer)
		if err != nil {
			dlog.Error(err)
		}
	} else {
		cb.SetLayer(layer)
		_, err := Draw(cb, stackLayer)
		if err != nil {
			dlog.Error(err)
		}
	}
	return cb
}

// DrawForTime draws and after d undraws an element
func DrawForTime(r Renderable, l int, d time.Duration) error {
	_, err := Draw(r, l)
	if err != nil {
		return err
	}
	go func(r Renderable, d time.Duration) {
		timing.DoAfter(d, func() {
			r.UnDraw()
		})
	}(r, d)
	return nil
}
