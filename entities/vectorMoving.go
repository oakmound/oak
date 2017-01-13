package entities

import (
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/physics"
)

type VectorMoving struct {
	Solid
	vMoving
}

func (vm *VectorMoving) Init() event.CID {
	cID := event.NextID(vm)
	vm.CID = cID
	return cID
}

func (vm *VectorMoving) ShiftVector(v *physics.Vector) {
	vm.Solid.ShiftPos(v.X, v.Y)
}

func (vm *VectorMoving) ApplyFriction(outsideFriction float64) {
	//Absolute friction is 1
	frictionScaler := 1 - (vm.Friction * outsideFriction)
	if frictionScaler > 1 {
		frictionScaler = 1
	} else if frictionScaler < 0 {
		frictionScaler = 0
	}
	vm.Delta.Scale(frictionScaler)
	if vm.Delta.Magnitude() < .01 {
		vm.Delta.Zero()
	}
}

type vMoving struct {
	Delta    *physics.Vector
	Speed    *physics.Vector
	Friction float64
}
