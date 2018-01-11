package entities

import (
	"strconv"

	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

// A Doodad is an entity composed of a position, a renderable, and a CallerID.
type Doodad struct {
	Point
	event.CID
	R render.Renderable
}

// NewDoodad returns a new doodad that is not drawn but is initialized.
// Passing a CID of 0 will initialize the entity as a Doodad. Passing
// any other CID will assume that the struct containing this doodad has
// already been initialized to the passed in CID.
// This applies to ALL NewX functions in entities which take in a CID.
func NewDoodad(x, y float64, r render.Renderable, CID event.CID) Doodad {
	if r != nil {
		r.SetPos(x, y)
	}
	d := Doodad{}
	d.Point = NewPoint(x, y)
	d.R = r
	d.CID = CID.Parse(&d)
	return d
}

// Init satisfies event.Entity
func (d *Doodad) Init() event.CID {
	d.CID = event.NextID(d)
	return d.CID
}

// GetID returns this Doodad's CID
// Consider: are these getters needed?
func (d *Doodad) GetID() event.CID {
	return d.CID
}

// GetRenderable returns this Doodad's Renderable
func (d *Doodad) GetRenderable() render.Renderable {
	return d.R
}

// SetRenderable sets this Doodad's renderable, drawing it.
// Todo:this automatic drawing doesn't really work with our
// two tiers of draw layers
func (d *Doodad) SetRenderable(r render.Renderable) {
	if d.R != nil {
		d.R.Undraw()
	}
	d.R = r
	render.Draw(d.R, d.R.GetLayer())
}

// Destroy cleans up the events, renderable and
// entity mapping for this Doodad
func (d *Doodad) Destroy() {
	if d.R != nil {
		d.R.Undraw()
	}
	d.CID.UnbindAll()
	event.DestroyEntity(int(d.CID))
}

// Overwrites

// SetPos both Sets logical position and renderable position
// The need for this sort of function is lessened with the introduction
// of vector attachement.
func (d *Doodad) SetPos(x, y float64) {
	d.SetLogicPos(x, y)
	d.R.SetPos(x, y)
}

func (d *Doodad) String() string {
	s := "Doodad: \nP{ "
	s += d.Point.String()
	s += " }\nID:{ "
	s += strconv.Itoa(int(d.CID))
	s += " }"
	return s
}
