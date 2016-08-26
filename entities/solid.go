package entities

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/collision"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
)

type Solid struct {
	Doodad
	W, H  float64
	Space *collision.Space
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

// Overwrites

func (s *Solid) Init() event.CID {
	cID := event.NextID(s)
	s.CID = cID
	return cID
}

func (s *Solid) SetPos(x float64, y float64) {
	s.SetLogicPos(x, y)
	s.R.SetPos(x, y)

	if s.Space != nil {
		collision.UpdateSpace(s.X, s.Y, s.W, s.H, s.Space)
	}
}

func (s *Solid) Destroy() {
	s.R.UnDraw()
	collision.Remove(s.Space)
	s.CID.UnbindAll()
	event.DestroyEntity(int(s.CID))
}
