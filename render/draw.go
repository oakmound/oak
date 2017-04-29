package render

import (
	"image/color"

	"time"

	"bitbucket.org/oakmoundstudio/oak/timing"
)

var (
	resetDraw       bool
	EmptyRenderable = NewColorBox(1, 1, color.RGBA{0, 0, 0, 0})
)

// LoadSpriteAndDraw is shorthand for LoadSprite
// followed by Draw.
func LoadSpriteAndDraw(filename string, l int) (Renderable, error) {
	s := LoadSprite(filename)
	return Draw(s, l)
}

// DrawColor is equivalent to LoadSpriteAndDraw,
// but with colorboxes.
func DrawColor(c color.Color, x1, y1, x2, y2 float64, l int) {
	cb := NewColorBox(int(x2), int(y2), c)
	cb.ShiftX(x1)
	cb.ShiftY(y1)
	Draw(cb, l)
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
