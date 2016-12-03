package entities

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/physics"
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

type vMoving struct {
	Delta *physics.Vector
	Speed *physics.Vector
	//Friction
}
