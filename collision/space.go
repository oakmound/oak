package collision

import (
	"log"
	"strconv"

	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/physics"
	"github.com/Sythe2o0/rtreego"
)

// ID Types constant
const (
	NONE = iota
	CID
	PID
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
	// Type represents which ID space the above ID
	// corresponds to.
	Type int
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

func (s *Space) Contains(other *Space) bool {
	//You contain another space if it is fully inside your space
	//If you are the same size and location as the space you are checking then you both contain eachother
	if s.GetX() > other.GetX() || s.GetX()+s.GetW() < other.GetX()+other.GetW() ||
		s.GetY() > other.GetY() || s.GetY()+s.GetH() < other.GetY()+other.GetH() {
		return false
	}
	return true
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

func (s *Space) UpdateLabel(classtype int) {
	rt.Delete(s)
	s.Label = classtype
	rt.Insert(s)
}

func (s *Space) OverlapVector(other *Space) physics.Vector {
	xover, yover := s.Overlap(other)
	return physics.NewVector(-xover, -yover)
}

func (s *Space) SubtractRect(x2, y2, w2, h2 float64) []*Space {
	x1 := s.GetX()
	y1 := s.GetY()
	w1 := s.GetW()
	h1 := s.GetH()

	// Left, Top, Right, Bottom
	// X, Y, W, H
	rects := [4][4]float64{}

	rects[0][0] = x1
	rects[0][1] = y1
	rects[0][2] = x2
	rects[0][3] = h1

	rects[1][0] = x1
	rects[1][1] = y1
	rects[1][2] = w1
	rects[1][3] = y2

	rects[2][0] = x1 + x2 + w2
	rects[2][1] = y1
	rects[2][2] = w1 - (x2 + w2)
	rects[2][3] = h1

	rects[3][0] = x1
	rects[3][1] = y1 + y2 + h2
	rects[3][2] = w1
	rects[3][3] = h1 - (y2 + h2)

	spaces := make([]*Space, 0)

	for _, r := range rects {
		if r[2] > 0 && r[3] > 0 {
			spaces = append(spaces, NewFullSpace(r[0], r[1], r[2], r[3], s.Label, s.CID))
		}
	}

	return spaces
}

func (s *Space) String() string {
	return strconv.FormatFloat(s.GetX(), 'f', 2, 32) + "," +
		strconv.FormatFloat(s.GetY(), 'f', 2, 32) + "," +
		strconv.FormatFloat(s.GetW(), 'f', 2, 32) + "," +
		strconv.FormatFloat(s.GetH(), 'f', 2, 32)
}

func NewUnassignedSpace(x, y, w, h float64) *Space {
	rect := NewRect(x, y, w, h)
	return &Space{
		Location: rect,
		Type:     NONE,
	}
}

func NewSpace(x, y, w, h float64, cID event.CID) *Space {
	rect := NewRect(x, y, w, h)
	return &Space{
		rect,
		-1,
		cID,
		CID,
	}
}

func NewLabeledSpace(x, y, w, h float64, l int) *Space {
	rect := NewRect(x, y, w, h)
	return &Space{
		Location: rect,
		Label:    l,
		Type:     NONE,
	}
}

func NewFullSpace(x, y, w, h float64, l int, cID event.CID) *Space {
	rect := NewRect(x, y, w, h)
	return &Space{
		rect,
		l,
		cID,
		CID,
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
