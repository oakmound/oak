package entities

import (
	"strconv"

	"bitbucket.org/oakmoundstudio/oak/collision"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/render"
)

type Reactive struct {
	Doodad
	W, H   float64
	RSpace *collision.ReactiveSpace
}

func NewReactive(x, y, w, h float64, r render.Renderable, cid event.CID) Reactive {
	return Reactive{
		Doodad: NewDoodad(x, y, r, cid),
		W:      w,
		H:      h,
		RSpace: collision.NewEmptyReactiveSpace(collision.NewSpace(x, y, w, h, cid)),
	}
}

func (r *Reactive) SetDim(w, h float64) {
	r.SetLogicDim(w, h)
	r.RSpace.SetDim(w, h)
}

func (r *Reactive) GetLogicDim() (float64, float64) {
	return r.W, r.H
}

func (r *Reactive) SetLogicDim(w, h float64) {
	r.W = w
	r.H = h
}

func (r *Reactive) SetSpace(sp *collision.ReactiveSpace) {
	collision.Remove(r.RSpace.Space)
	r.RSpace = sp
	collision.Add(r.RSpace.Space)
}

func (r *Reactive) GetSpace() *collision.Space {
	return r.RSpace.Space
}

// Overwrites

func (r *Reactive) Init() event.CID {
	cID := event.NextID(r)
	r.CID = cID
	return cID
}

func (r *Reactive) ShiftPos(x, y float64) {
	r.SetPos(r.X()+x, r.Y()+y)
}

func (r *Reactive) SetPos(x, y float64) {
	r.SetLogicPos(x, y)
	r.R.SetPos(x, y)
	collision.UpdateSpace(r.X(), r.Y(), r.W, r.H, r.RSpace.Space)
}

func (r *Reactive) Destroy() {
	r.R.UnDraw()
	collision.Remove(r.RSpace.Space)
	r.CID.UnbindAll()
	event.DestroyEntity(int(r.CID))
}

func (r *Reactive) String() string {
	st := "Reactive:\n{"
	st += r.Doodad.String()
	st += " }, \n"
	w := strconv.FormatFloat(r.W, 'f', 2, 32)
	h := strconv.FormatFloat(r.H, 'f', 2, 32)
	st += "W: " + w + ", H: " + h
	st += ",\nS:{ "
	st += r.RSpace.String()
	st += "}"
	return st
}
