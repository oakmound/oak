package move

import (
	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/alg/floatgeom"
	"github.com/oakmound/oak/v2/key"
	"github.com/oakmound/oak/v2/physics"
)

// WASD moves the given mover based on its speed as W,A,S, and D are pressed
func WASD(mvr Mover) {
	TopDown(mvr, key.W, key.S, key.A, key.D)
}

// Arrows moves the given mover based on its speed as the arrow keys are pressed
func Arrows(mvr Mover) {
	TopDown(mvr, key.UpArrow, key.DownArrow, key.LeftArrow, key.RightAlt)
}

// TopDown moves the given mover based on its speed as the given keys are pressed
func TopDown(mvr Mover, up, down, left, right string) {
	delta := mvr.GetDelta()
	vec := mvr.Vec()
	spd := mvr.GetSpeed()

	delta.Zero()
	if oak.IsDown(up) {
		delta.Add(physics.NewVector(0, -spd.Y()))
	}
	if oak.IsDown(down) {
		delta.Add(physics.NewVector(0, spd.Y()))
	}
	if oak.IsDown(left) {
		delta.Add(physics.NewVector(-spd.X(), 0))
	}
	if oak.IsDown(right) {
		delta.Add(physics.NewVector(spd.X(), 0))
	}
	vec.Add(delta)
	mvr.GetRenderable().SetPos(vec.X(), vec.Y())
	sp := mvr.GetSpace()
	sp.Update(vec.X(), vec.Y(), sp.GetW(), sp.GetH())
}

// CenterScreenOn will cause the screen to center on the given mover, obeying
// viewport limits if they have been set previously
func CenterScreenOn(mvr Mover) {
	vec := mvr.Vec()
	oak.SetScreen(
		int(vec.X())-oak.ScreenWidth/2,
		int(vec.Y())-oak.ScreenHeight/2,
	)
}

// Limit restricts the movement of the mover to stay within a given rectangle
func Limit(mvr Mover, rect floatgeom.Rect2) {
	vec := mvr.Vec()
	w, h := mvr.GetRenderable().GetDims()
	wf := float64(w)
	hf := float64(h)
	if vec.X() < rect.Min.X() {
		vec.SetX(rect.Min.X())
	} else if vec.X() > rect.Max.X()-wf {
		vec.SetX(rect.Max.X() - wf)
	}
	if vec.Y() < rect.Min.Y() {
		vec.SetY(rect.Min.Y())
	} else if vec.Y() > rect.Max.Y()-hf {
		vec.SetY(rect.Max.Y() - hf)
	}
}
