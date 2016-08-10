package entities

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/collision"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
)

type Point struct {
	X, Y float64
}

func (p *Point) SetPos(x, y float64) {
	p.SetLogicPos(x, y)
}
func (p *Point) GetLogicPos() (float64, float64) {
	return p.X, p.Y
}
func (p *Point) SetLogicPos(x, y float64) {
	p.X = x
	p.Y = y
}

////////////////////////

type Doodad struct {
	Point
	R   render.Renderable
	CID event.CID
}

func (d *Doodad) Init() event.CID {
	cID := event.NextID(d)
	d.CID = cID
	return cID
}

func (d *Doodad) GetRenderable() render.Renderable {
	return d.R
}

func (d *Doodad) GetID() event.CID {
	return d.CID
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

////////////////

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
