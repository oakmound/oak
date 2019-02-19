package move

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/physics"
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
	mvr.GetSpace().Update(vec.X(), vec.Y(), 16, 16)
}
