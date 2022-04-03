package entities

import (
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/physics"
	"github.com/oakmound/oak/v3/render"
)

// A Moving is a Solid that also keeps track of a speed and a delta vector
type Moving struct {
	Solid
	vMoving
}

// NewMoving returns a new Moving
func NewMoving(x, y, w, h float64, r render.Renderable, tree *collision.Tree, cid event.CallerID, friction float64) *Moving {
	m := &Moving{}
	if cid == 0 {
		m.CallerID = event.DefaultCallerMap.Register(m)
	} else {
		m.CallerID = cid
	}
	m.Solid = *NewSolid(x, y, w, h, r, tree, m.CallerID)
	m.vMoving = vMoving{
		Delta:    physics.NewVector(0, 0),
		Speed:    physics.NewVector(0, 0),
		Friction: friction,
	}
	return m
}

func (m *Moving) CID() event.CallerID {
	return m.CallerID
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

// GetSpeed returns this moving's speed
func (v vMoving) GetSpeed() physics.Vector {
	return v.Speed
}
