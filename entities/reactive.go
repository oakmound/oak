package entities

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/collision"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
)

type Interactive struct {
	Reactive
	moving
}

type Reactive struct {
	Doodad
	W, H   float64
	RSpace *collision.ReactiveSpace
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
	collision.Remove(r.RSpace.Space())
	r.RSpace = sp
	collision.Add(r.RSpace.Space())
}

func (r *Reactive) GetSpace() *collision.ReactiveSpace {
	return r.RSpace
}

// Overwrites

func (r *Reactive) SetPos(x float64, y float64) {
	r.SetLogicPos(x, y)
	r.R.SetPos(x, y)

	if r.RSpace != nil {
		collision.UpdateSpace(r.X, r.Y, r.W, r.H, r.RSpace.Space())
	}
}

func (r *Reactive) Destroy() {
	r.R.UnDraw()
	collision.Remove(r.RSpace.Space())
	r.CID.UnbindAll()
	event.DestroyEntity(int(r.CID))
}
