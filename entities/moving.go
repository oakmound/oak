package entities

import (
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
)

// A Moving is a Solid that also keeps track of a speed and a delta vector
type Moving struct {
	Solid
	vMoving
}

// NewMoving returns a new Moving
func NewMoving(x, y, w, h float64, r render.Renderable, cid event.CID, friction float64) Moving {
	m := Moving{}
	cid = cid.Parse(&m)
	m.Solid = NewSolid(x, y, w, h, r, cid)
	m.vMoving = vMoving{
		Delta:    physics.NewVector(0, 0),
		Speed:    physics.NewVector(0, 0),
		Friction: friction,
	}
	return m
}

// Init satisfies event.Entity
func (m *Moving) Init() event.CID {
	m.CID = event.NextID(m)
	return m.CID
}

// ShiftVector probably shouldn't be on moving but it lets you
// ShiftPos by a given vector
func (m *Moving) ShiftVector(v physics.Vector) {
	m.Solid.ShiftPos(v.X(), v.Y())
}

// ApplyFriction modifies a moving's delta by combining
// environmental friction with the moving's base friction
// and scaling down the delta by the combined result.
func (v *vMoving) ApplyFriction(outsideFriction float64) {
	//Absolute friction is 1
	frictionScaler := 1 - (v.Friction * outsideFriction)
	if frictionScaler > 1 {
		frictionScaler = 1
	} else if frictionScaler < 0 {
		frictionScaler = 0
	}
	v.Delta.Scale(frictionScaler)
	if v.Delta.Magnitude() < .01 {
		v.Delta.Zero()
	}
}

type vMoving struct {
	Delta    physics.Vector
	Speed    physics.Vector
	Friction float64
}

// GetDelta returns this moving's delta
func (v vMoving) GetDelta() physics.Vector {
	return v.Delta
}
