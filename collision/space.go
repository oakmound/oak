package collision

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	//"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
	"github.com/dhconnelly/rtreego"
	//"image/color"
)

// Spaces are a rectangle
// with a couple of ways of identifying
// the underlying object.
type Space struct {
	Location *rtreego.Rect
	// A label can store type information.
	// Recommended to use with an enum.
	Label int
	// A CID can be used to get the exact
	// entity which this rectangle belongs to.
	CID event.CID
}

// Bounds satisfies the rtreego.Spatial interface.
func (s *Space) Bounds() *rtreego.Rect {
	return s.Location
}

func (s *Space) GetX() float64 {
	return s.Location.PointCoord(0)
}

func (s *Space) GetY() float64 {
	return s.Location.PointCoord(1)
}

func (s *Space) GetW() float64 {
	return s.Location.LengthsCoord(0)
}

func (s *Space) GetH() float64 {
	return s.Location.LengthsCoord(1)
}

func (s *Space) GetCenter() (float64, float64) {
	return s.GetX() + s.GetW()/2, s.GetY() + s.GetH()/2
}

func (s *Space) GetPos() (float64, float64) {
	return s.Location.PointCoord(1), s.Location.PointCoord(0)
}

func (s *Space) Above(other *Space) float64 {
	return other.GetY() - s.GetY()
}

func (s *Space) Below(other *Space) float64 {
	return s.GetY() - other.GetY()
}

func (s *Space) LeftOf(other *Space) float64 {
	return other.GetX() - s.GetX()
}

func (s *Space) RightOf(other *Space) float64 {
	return s.GetX() - other.GetX()
}

func (s *Space) SetDim(w, h float64) {
	s.Update(s.GetX(), s.GetY(), w, h)
}

func (s *Space) Update(x, y, w, h float64) {
	loc := NewRect(x, y, w, h)
	rt.Delete(s)
	s.Location = loc
	rt.Insert(s)
}

func NewUnassignedSpace(x, y, w, h float64) *Space {
	//render.DrawColor(color.RGBA{128, 0, 128, 100}, x, y, w, h, 10)
	rect := NewRect(x, y, w, h)
	return &Space{Location: rect}
}

func NewSpace(x, y, w, h float64, cID event.CID) *Space {
	rect := NewRect(x, y, w, h)
	return &Space{
		rect,
		-1,
		cID,
	}
}

func NewLabeledSpace(x, y, w, h float64, l int) *Space {
	rect := NewRect(x, y, w, h)
	return &Space{
		Location: rect,
		Label:    l,
	}
}

func NewFullSpace(x, y, w, h float64, l int, cID event.CID) *Space {
	rect := NewRect(x, y, w, h)
	return &Space{
		rect,
		l,
		cID,
	}
}
