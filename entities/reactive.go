package entities

import (
	"strconv"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

// Reactive is parallel to Solid, but has a Reactive collision space instead of
// a regular collision space
type Reactive struct {
	Doodad
	W, H   float64
	RSpace *collision.ReactiveSpace
	Tree   *collision.Tree
}

// NewReactive returns a new Reactive struct. The added space will
// be added to the input tree, or DefTree if none is given.
func NewReactive(x, y, w, h float64, r render.Renderable, tree *collision.Tree, cid event.CID) Reactive {
	rct := Reactive{}
	cid = cid.Parse(&rct)
	rct.Doodad = NewDoodad(x, y, r, cid)
	rct.W = w
	rct.H = h
	rct.RSpace = collision.NewEmptyReactiveSpace(collision.NewSpace(x, y, w, h, cid))
	if tree == nil {
		tree = collision.DefTree
	}
	rct.Tree = tree
	rct.Tree.Add(rct.RSpace.Space)
	return rct
}

// SetDim sets the dimensions of this reactive's space and it's logical dimensions
func (r *Reactive) SetDim(w, h float64) {
	r.SetLogicDim(w, h)
	r.RSpace.SetDim(w, h)
}

// GetLogicDim returns this Reactive's width and height
// todo: move wh into their own struct to compose into solid and reactive
func (r *Reactive) GetLogicDim() (float64, float64) {
	return r.W, r.H
}

// SetLogicDim sets the logical width and height of this reactive
// without changing the real dimensions of its collision space
func (r *Reactive) SetLogicDim(w, h float64) {
	r.W = w
	r.H = h
}

// SetSpace sets this reactive's collision space to the given reactive space,
// updating it's collision tree to include it.
func (r *Reactive) SetSpace(sp *collision.ReactiveSpace) {
	r.Tree.Remove(r.RSpace.Space)
	r.RSpace = sp
	r.Tree.Add(r.RSpace.Space)
}

// GetSpace returns this reactive's space underlying its RSpace
func (r *Reactive) GetSpace() *collision.Space {
	return r.RSpace.Space
}

// Overwrites

// Init satisfies event.Entity
func (r *Reactive) Init() event.CID {
	r.CID = event.NextID(r)
	return r.CID
}

// ShiftPos acts like SetPos if given r.X()+x, r.Y()+y
func (r *Reactive) ShiftPos(x, y float64) {
	r.SetPos(r.X()+x, r.Y()+y)
}

// SetPos sets this reactive's logical, renderable, and collision position to be x,y
func (r *Reactive) SetPos(x, y float64) {
	r.SetLogicPos(x, y)
	r.R.SetPos(x, y)
	r.Tree.UpdateSpace(r.X(), r.Y(), r.W, r.H, r.RSpace.Space)
}

// Destroy destroys this reactive's doodad component and removes its space
// from it's collision tree
func (r *Reactive) Destroy() {
	r.Doodad.Destroy()
	r.Tree.Remove(r.RSpace.Space)
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
