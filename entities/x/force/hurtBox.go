package force

import (
	"image/color"
	"time"

	"github.com/oakmound/lowrez17/game/layers"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/timing"
)

type hurtBox struct {
	*DirectionSpace
}

// NewHurtBox creates a temporary collision space with a given force it should
// apply to objects it collides with
func NewHurtBox(x, y, w, h float64, duration time.Duration, l collision.Label, fv physics.ForceVector) {
	hb := new(hurtBox)
	hb.DirectionSpace = NewDirectionSpace(collision.NewLabeledSpace(x, y, w, h, l), fv)
	collision.Add(hb.Space)
	go timing.DoAfter(duration, func() {
		collision.Remove(hb.Space)
	})
}

// NewHurtColor creates a temporary collision space with a given force it should
// apply to objects it collides with. The box is rendered as the given color
func NewHurtColor(x, y, w, h float64, duration time.Duration, l collision.Label, fv physics.ForceVector, c color.Color) {
	cb := render.NewColorBox(int(w), int(h), c)
	NewHurtDisplay(x, y, w, h, duration, l, fv, cb)
}

// NewHurtDisplay creates a temporary collision space with a given force it should
// apply to objects it collides with. The box is rendered as the given renderable.
// The input renderable is not copied before it is drawn.
func NewHurtDisplay(x, y, w, h float64, duration time.Duration, l collision.Label, fv physics.ForceVector, r render.Renderable) {
	hb := new(hurtBox)
	hb.DirectionSpace = NewDirectionSpace(collision.NewLabeledSpace(x, y, w, h, l), fv)
	collision.Add(hb.Space)
	r.SetPos(x, y)
	render.Draw(r, layers.DebugLayer)
	go timing.DoAfter(duration, func() {
		collision.Remove(hb.Space)
		r.Undraw()
	})
}
