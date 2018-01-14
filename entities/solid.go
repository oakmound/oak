package entities

import (
	"strconv"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

// A Solid is a Doodad with a width, height, and collision space.
type Solid struct {
	Doodad
	W, H  float64
	Space *collision.Space
	Tree  *collision.Tree
}

// NewSolid returns an initialized Solid that is not drawn and whose space
// belongs to the given collision tree. If nil is given as the tree, it will
// belong to collision.DefTree
func NewSolid(x, y, w, h float64, r render.Renderable, tree *collision.Tree, cid event.CID) Solid {
	s := Solid{}
	cid = cid.Parse(&s)
	s.Doodad = NewDoodad(x, y, r, cid)
	s.W = w
	s.H = h
	if tree == nil {
		tree = collision.DefTree
	}
	s.Tree = tree
	s.Space = collision.NewSpace(x, y, w, h, cid)
	s.Tree.Add(s.Space)
	return s
}

// SetDim sets the logical dimensions of the solid and the real
// dimensions on the solid's space
func (s *Solid) SetDim(w, h float64) {
	s.SetLogicDim(w, h)
	s.Space.SetDim(w, h)
}

// GetLogicDim will return the width and height of the Solid
func (s *Solid) GetLogicDim() (float64, float64) {
	return s.W, s.H
}

// SetLogicDim sets the width and height of the solid
func (s *Solid) SetLogicDim(w, h float64) {
	s.W = w
	s.H = h
}

// SetSpace assigns a solid a collision space and puts it in this Solid's Tree
func (s *Solid) SetSpace(sp *collision.Space) {
	s.Tree.Remove(s.Space)
	s.Space = sp
	s.Tree.Add(s.Space)
}

// GetSpace returns a solid's collision space
func (s *Solid) GetSpace() *collision.Space {
	return s.Space
}

// ShiftX moves a solid by x along the x axis
func (s *Solid) ShiftX(x float64) {
	s.SetPos(s.X()+x, s.Y())
}

// ShiftY moves a solid by y along the y axis
func (s *Solid) ShiftY(y float64) {
	s.SetPos(s.X(), s.Y()+y)
}

// ShiftPos moves a solid by (x,y)
func (s *Solid) ShiftPos(x, y float64) {
	s.SetPos(s.X()+x, s.Y()+y)
}

// Overwrites

// Init satisfies event.Entity
func (s *Solid) Init() event.CID {
	s.CID = event.NextID(s)
	return s.CID
}

// SetPos sets the position of the collision space, the logical position,
// and the renderable position of the solid.
func (s *Solid) SetPos(x float64, y float64) {
	s.SetLogicPos(x, y)
	if s.R != nil {
		s.R.SetPos(x, y)
	}
	s.Tree.UpdateSpace(s.X(), s.Y(), s.W, s.H, s.Space)
}

// Destroy removes this solid's collision space from it's Tree
// and destroys the doodad portion of the solid.
func (s *Solid) Destroy() {
	s.Doodad.Destroy()
	s.Tree.Remove(s.Space)
}

func (s *Solid) String() string {
	st := "Solid:\n{"
	st += s.Doodad.String()
	st += "},\n"
	w := strconv.FormatFloat(s.W, 'f', 2, 32)
	h := strconv.FormatFloat(s.H, 'f', 2, 32)
	st += "W: " + w + ", H: " + h
	st += ",\nS:{"
	st += s.Space.String()
	st += "}"
	return st
}
