package entities

import (
	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/key"
)

// WASD moves the given mover based on its speed as W,A,S, and D are pressed
func WASD(mvr *Entity) {
	TopDown(mvr, key.W, key.S, key.A, key.D)
}

// Arrows moves the given mover based on its speed as the arrow keys are pressed
func Arrows(mvr *Entity) {
	TopDown(mvr, key.UpArrow, key.DownArrow, key.LeftArrow, key.RightAlt)
}

// TopDown moves the given mover based on its speed as the given keys are pressed
func TopDown(mvr *Entity, up, down, left, right key.Code) {
	mvr.Delta = floatgeom.Point2{}
	if mvr.ctx.KeyState.IsDown(up) {
		mvr.Delta[1] -= mvr.Speed[1]
	}
	if mvr.ctx.KeyState.IsDown(down) {
		mvr.Delta[1] += mvr.Speed[1]
	}
	if mvr.ctx.KeyState.IsDown(left) {
		mvr.Delta[0] -= mvr.Speed[0]
	}
	if mvr.ctx.KeyState.IsDown(right) {
		mvr.Delta[0] += mvr.Speed[0]
	}
	mvr.ShiftDelta()
}

// CenterScreenOn will cause the screen to center on the given mover, obeying
// viewport limits if they have been set previously
func CenterScreenOn(mvr *Entity) {
	bds := mvr.ctx.Window.Bounds()
	pos := intgeom.Point2{int(mvr.X()), int(mvr.Y())}
	target := pos.Sub(bds).DivConst(2)
	mvr.ctx.Window.SetViewport(target)
}

// Limit restricts the movement of the mover to stay within a given rectangle
func Limit(mvr *Entity, rect floatgeom.Rect2) {
	wf := mvr.W()
	hf := mvr.H()
	if mvr.X() < rect.Min.X() {
		mvr.SetX(rect.Min.X())
	} else if mvr.X() > rect.Max.X()-wf {
		mvr.SetX(rect.Max.X() - wf)
	}
	if mvr.Y() < rect.Min.Y() {
		mvr.SetY(rect.Min.Y())
	} else if mvr.Y() > rect.Max.Y()-hf {
		mvr.SetY(rect.Max.Y() - hf)
	}
}
