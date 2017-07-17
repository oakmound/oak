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
}

// NewSolid returns an initialized Solid that is not drawn and whose space
// belongs to no collision tree
func NewSolid(x, y, w, h float64, r render.Renderable, cid event.CID) Solid {
	s := Solid{}
	// Todo: use this parse structure on everything else
	cid = cid.Parse(&s)
	s.Doodad = NewDoodad(x, y, r, cid)
	s.W = w
	s.H = h
	s.Space = collision.NewSpace(x, y, w, h, cid)
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

// SetSpace assigns a solid a collision space and puts it in the default
// collision tree. This is legacy behavior that should change.
func (s *Solid) SetSpace(sp *collision.Space) {
	// Todo: these functions should not be here, or this should know what
	// tree it belongs to!
	collision.Remove(s.Space)
	s.Space = sp
	collision.Add(s.Space)
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
	// This uses the legacy collision tree behavior
	collision.UpdateSpace(s.X(), s.Y(), s.W, s.H, s.Space)
}

// Destroy removes this solid's collision space from the default tree (todo)
// and destroys the doodad portion of the solid.
func (s *Solid) Destroy() {
	s.Doodad.Destroy()
	collision.Remove(s.Space)
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
