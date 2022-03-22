package entities

import (
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/render"
)

// A Doodad is an entity composed of a position, a renderable, and a CallerID.
type Doodad struct {
	Point
	event.CallerID
	R render.Renderable
}

// NewDoodad returns a new doodad that is not drawn but is initialized.
// Passing a CID of 0 will initialize the entity as a Doodad. Passing
// any other CID will assume that the struct containing this doodad has
// already been initialized to the passed in CID.
// This applies to ALL NewX functions in entities which take in a CID.
func NewDoodad(x, y float64, r render.Renderable, cid event.CallerID) *Doodad {
	if r != nil {
		r.SetPos(x, y)
	}
	d := &Doodad{}
	d.Point = *NewPoint(x, y)
	d.R = r
	if cid == 0 {
		d.CallerID = event.DefaultCallerMap.Register(d)
	} else {
		d.CallerID = cid
	}
	return d
}

func (d *Doodad) CID() event.CallerID {
	return d.CallerID
}

// Destroy cleans up the events, renderable and
// entity mapping for this Doodad
func (d *Doodad) Destroy() {
	if d.R != nil {
		d.R.Undraw()
	}
	event.DefaultBus.UnbindAllFrom(d.CallerID)
	event.DefaultCallerMap.DestroyEntity(d.CallerID)
}

// Overwrites

// SetPos both Sets logical position and renderable position
// The need for this sort of function is lessened with the introduction
// of vector attachement.
func (d *Doodad) SetPos(x, y float64) {
	d.SetLogicPos(x, y)
	if d.R != nil {
		d.R.SetPos(x, y)
	}
}
