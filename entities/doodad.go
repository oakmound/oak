package entities

import (
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/render"
)

type Doodad struct {
	Point
	event.CID
	R render.Renderable
}

func NewDoodad(x, y float64, r render.Renderable, CID event.CID) Doodad {
	if r != nil {
		r.SetPos(x, y)
	}
	return Doodad{
		Point: NewPoint(x, y),
		R:     r,
		CID:   CID,
	}
}

func (d *Doodad) Init() event.CID {
	cID := event.NextID(d)
	d.CID = cID
	return cID
}

func (d *Doodad) GetID() event.CID {
	return d.CID
}

func (d *Doodad) GetRenderable() render.Renderable {
	return d.R
}

func (d *Doodad) SetRenderable(r render.Renderable) {
	if d.R != nil {
		d.R.UnDraw()
	}
	d.R = r
	render.Draw(d.R, d.R.GetLayer())
}

func (d *Doodad) Destroy() {
	d.R.UnDraw()
	d.CID.UnbindAll()
	event.DestroyEntity(int(d.CID))
}

// Overwrites
func (d *Doodad) SetPos(x, y float64) {
	d.SetLogicPos(x, y)
	d.R.SetPos(x, y)
}

func (d *Doodad) String() string {
	s := "Doodad: \nP{ "
	s += d.Point.String()
	s += " }\nR:{ "
	s += d.R.String()
	s += " }\nID:{ "
	s += d.CID.String()
	s += " }"
	return s
}
