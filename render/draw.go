package render

import (
	"image/color"

	"time"

	"bitbucket.org/oakmoundstudio/oak/timing"
)

var (
	resetDraw bool
	// EmptyRenderable is a simple renderable that can be used
	// for pseudo-nil renderables that need to be something
	emptyRenderable = NewColorBox(1, 1, color.RGBA{0, 0, 0, 0})
)

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
		Draw(cb, layer)
	} else {
		cb.SetLayer(layer)
		Draw(cb, stackLayer)
	}
	return cb
}

// UndrawAfter will trigger a renderable's undraw function
// after a given time has passed
func UndrawAfter(r Renderable, t time.Duration) {
	go func(r Renderable, t time.Duration) {
		timing.DoAfter(t, func() {
			r.UnDraw()
		})
	}(r, t)
}

// DrawForTime is a wrapper for Draw and UndrawAfter
func DrawForTime(r Renderable, l int, t time.Duration) {
	Draw(r, l)
	UndrawAfter(r, t) 
}