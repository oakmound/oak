package entities

import (
	"strconv"

	"bitbucket.org/oakmoundstudio/oak/collision"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/render"
)

type Solid struct {
	Doodad
	W, H  float64
	Space *collision.Space
}

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

func (s *Solid) SetDim(w, h float64) {
	s.SetLogicDim(w, h)
	s.Space.SetDim(w, h)
}

func (s *Solid) GetLogicDim() (float64, float64) {
	return s.W, s.H
}

func (s *Solid) SetLogicDim(w, h float64) {
	s.W = w
	s.H = h
}

func (s *Solid) SetSpace(sp *collision.Space) {
	collision.Remove(s.Space)
	s.Space = sp
	collision.Add(s.Space)
}

func (s *Solid) GetSpace() *collision.Space {
	return s.Space
}

func (s *Solid) ShiftX(x float64) {
	s.SetPos(s.X()+x, s.Y())
}

func (s *Solid) ShiftY(y float64) {
	s.SetPos(s.X(), s.Y()+y)
}

func (s *Solid) ShiftPos(x, y float64) {
	s.SetPos(s.X()+x, s.Y()+y)
}

// Overwrites

func (s *Solid) Init() event.CID {
	cID := event.NextID(s)
	s.CID = cID
	return cID
}

func (s *Solid) SetPos(x float64, y float64) {
	s.SetLogicPos(x, y)
	if s.R != nil {
		s.R.SetPos(x, y)
	}
	collision.UpdateSpace(s.X(), s.Y(), s.W, s.H, s.Space)
}

func (s *Solid) Destroy() {
	if s.R != nil {
		s.R.UnDraw()
	}
	collision.Remove(s.Space)
	s.CID.UnbindAll()
	event.DestroyEntity(int(s.CID))
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
