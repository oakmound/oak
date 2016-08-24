package entities

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/collision"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
)

type Point struct {
	X, Y float64
}

func (p *Point) GetX() float64 {
	return p.X
}
func (p *Point) GetY() float64 {
	return p.Y
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

func (d *Doodad) GetID() event.CID {
	return d.CID
}

func (d *Doodad) GetRenderable() render.Renderable {
	return d.R
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

//////////////

type moving struct {
	DX, DY, SpeedX, SpeedY float64
}

func (m *moving) GetDX() float64 {
	return m.DX
}
func (m *moving) GetDY() float64 {
	return m.DY
}
func (m *moving) SetDXY(x, y float64) {
	m.DX = x
	m.DY = y
}
func (m *moving) GetSpeedX() float64 {
	return m.SpeedX
}
func (m *moving) GetSpeedY() float64 {
	return m.SpeedY
}
func (m *moving) SetSpeedXY(x, y float64) {
	m.SpeedX = x
	m.SpeedY = y
}

type Moving struct {
	Solid
	moving
}

//////////

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

/////////

type Interactive struct {
	Reactive
	moving
}
