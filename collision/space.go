package collision

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	//"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
	"github.com/Sythe2o0/rtreego"
	//"image/color"
	"log"
	"strconv"
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

func (s *Space) Overlap(other *Space) (xOver, yOver float64) {
	if s.GetX() > other.GetX() {
		x2 := other.GetX() + other.GetW()
		if s.GetX() < x2 {
			xOver = s.GetX() - x2
		}
	} else {
		x2 := s.GetX() + s.GetW()
		if other.GetX() < x2 {
			xOver = x2 - other.GetX()
		}
	}
	if s.GetY() > other.GetY() {
		y2 := other.GetY() + other.GetH()
		if s.GetY() < y2 {
			yOver = s.GetY() - y2
		}
	} else {
		y2 := s.GetY() + s.GetW()
		if other.GetY() < y2 {
			yOver = y2 - other.GetY()
		}
	}
	return
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

func (s *Space) String() string {
	return strconv.FormatFloat(s.GetX(), 'f', 2, 32) + "," +
		strconv.FormatFloat(s.GetY(), 'f', 2, 32) + "," +
		strconv.FormatFloat(s.GetW(), 'f', 2, 32) + "," +
		strconv.FormatFloat(s.GetH(), 'f', 2, 32)
}

func NewUnassignedSpace(x, y, w, h float64) *Space {
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

// NewRect is a wrapper around rtreego.NewRect,
// casting the given x,y to an rtreego.Point.
// Used to not expose rtreego.Point to the user.
func NewRect(x, y, w, h float64) *rtreego.Rect {
	rect, err := rtreego.NewRect(rtreego.Point{x, y}, [3]float64{w, h, 1})
	if err != nil {
		log.Fatal(err)
	}
	return &rect
}
