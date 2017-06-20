package entities

import (
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/physics"
	"bitbucket.org/oakmoundstudio/oak/render"
)

type Moving struct {
	Solid
	vMoving
}

func NewMoving(x, y, w, h float64, r render.Renderable, cid event.CID, friction float64) Moving {
	return Moving{
		Solid: NewSolid(x, y, w, h, r, cid),
		vMoving: vMoving{
			Delta:    physics.NewVector(0, 0),
			Speed:    physics.NewVector(0, 0),
			Friction: friction,
		},
	}
}

func (m *Moving) Init() event.CID {
	cID := event.NextID(m)
	m.CID = cID
	return cID
}

func (m *Moving) ShiftVector(v physics.Vector) {
	m.Solid.ShiftPos(v.X(), v.Y())
}

func (m *Moving) ApplyFriction(outsideFriction float64) {
	//Absolute friction is 1
	frictionScaler := 1 - (m.Friction * outsideFriction)
	if frictionScaler > 1 {
		frictionScaler = 1
	} else if frictionScaler < 0 {
		frictionScaler = 0
	}
	m.Delta.Scale(frictionScaler)
	if m.Delta.Magnitude() < .01 {
		m.Delta.Zero()
	}
}

type vMoving struct {
	Delta    physics.Vector
	Speed    physics.Vector
	Friction float64
}

func (v vMoving) GetDelta() physics.Vector {
	return v.Delta
}
