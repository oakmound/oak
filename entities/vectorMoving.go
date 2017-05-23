package entities

import (
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/physics"
	"bitbucket.org/oakmoundstudio/oak/render"
)

type VectorMoving struct {
	Solid
	vMoving
}

func NewVectorMoving(x, y, w, h float64, r render.Renderable, cid event.CID, friction float64) VectorMoving {
	return VectorMoving{
		Solid: NewSolid(x, y, w, h, r, cid),
		vMoving: vMoving{
			Delta:    physics.NewVector(0, 0),
			Speed:    physics.NewVector(0, 0),
			Friction: friction,
		},
	}
}

func (vm *VectorMoving) Init() event.CID {
	cID := event.NextID(vm)
	vm.CID = cID
	return cID
}

func (vm *VectorMoving) ShiftVector(v physics.Vector) {
	vm.Solid.ShiftPos(v.X(), v.Y())
}

func (vm *VectorMoving) ApplyFriction(outsideFriction float64) {
	//Absolute friction is 1
	frictionScaler := 1 - (vm.Friction * outsideFriction)
	if frictionScaler > 1 {
		frictionScaler = 1
	} else if frictionScaler < 0 {
		frictionScaler = 0
	}
	vm.Delta = vm.Delta.Scale(frictionScaler)
	if vm.Delta.Magnitude() < .01 {
		vm.Delta = vm.Delta.Zero()
	}
}

type vMoving struct {
	Delta    physics.Vector
	Speed    physics.Vector
	Friction float64
}
