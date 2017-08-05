package collision

import (
	"fmt"
	"strconv"

	"github.com/Sythe2o0/rtreego"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
)

// ID Types constant
const (
	NONE = iota
	CID
	PID
)

// A Space is a rectangle
// with a couple of ways of identifying
// an underlying object.
type Space struct {
	Location *rtreego.Rect
	// A label can store type information.
	// Recommended to use with an enum.
	Label Label
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

// GetX returns a space's x position (leftmost)
func (s *Space) GetX() float64 {
	return s.Location.PointCoord(0)
}

// GetY returns a space's y position (upmost)
func (s *Space) GetY() float64 {
	return s.Location.PointCoord(1)
}

// GetW returns a space's width (rightmost x - leftmost x)
func (s *Space) GetW() float64 {
	return s.Location.LengthsCoord(0)
}

// GetH returns a space's height (upper y - lower y)
func (s *Space) GetH() float64 {
	return s.Location.LengthsCoord(1)
}

// GetCenter returns the center point of the space
func (s *Space) GetCenter() (float64, float64) {
	return s.GetX() + s.GetW()/2, s.GetY() + s.GetH()/2
}

// GetPos returns both y and x
func (s *Space) GetPos() (float64, float64) {
	return s.Location.PointCoord(1), s.Location.PointCoord(0)
}

// Above returns how much above this space another space is
// Important note: (10,10) is Above (10,20), because in oak's
// display, lower y values are higher than higher y values.
func (s *Space) Above(other *Space) float64 {
	return other.GetY() - s.GetY()
}

// Below returns how much below this space another space is,
// Equivalent to -1 * Above
func (s *Space) Below(other *Space) float64 {
	return s.GetY() - other.GetY()
}

// Contains returns whether this space contains other
func (s *Space) Contains(other *Space) bool {
	//You contain another space if it is fully inside your space
	//If you are the same size and location as the space you are checking then you both contain eachother
	if s.GetX() > other.GetX() || s.GetX()+s.GetW() < other.GetX()+other.GetW() ||
		s.GetY() > other.GetY() || s.GetY()+s.GetH() < other.GetY()+other.GetH() {
		return false
	}
	return true
}

// LeftOf returns how far to the left other is of this space
func (s *Space) LeftOf(other *Space) float64 {
	return other.GetX() - s.GetX()
}

// RightOf returns how far to the right other is of this space.
// Equivalent to -1 * LeftOf
func (s *Space) RightOf(other *Space) float64 {
	return s.GetX() - other.GetX()
}

// Overlap returns how much this space overlaps with another space
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

// OverlapVector returns Overlap as a vector
func (s *Space) OverlapVector(other *Space) physics.Vector {
	xover, yover := s.Overlap(other)
	// Todo: why are we multiplying by -1 here, shouldn't that
	// also be happening in Overlap at least?
	return physics.NewVector(-xover, -yover)
}

// SubtractRect removes a rectangle from this rectangle and
// returns the rectangles remaining after the portion has been
// removed. The input x,y is relative to the original space:
// Example: removing 1,1 from 10,10 -> 12,12 is OK, but removing
// 11,11 from 10,10 -> 12,12 will not act as expected.
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

	// Todo: these spaces overlap on the corners. We could remove that.
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

	var spaces []*Space

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
		strconv.FormatFloat(s.GetH(), 'f', 2, 32) + "::" +
		strconv.Itoa(int(s.CID)) + "::" + fmt.Sprintf("%p", s)
}

// NewUnassignedSpace returns a space that just has a rectangle
func NewUnassignedSpace(x, y, w, h float64) *Space {
	return NewLabeledSpace(x, y, w, h, NilLabel)
}

// NewSpace returns a space with an associated caller id
func NewSpace(x, y, w, h float64, cID event.CID) *Space {
	return NewFullSpace(x, y, w, h, NilLabel, cID)
}

// NewLabeledSpace returns a space with an associated integer label
func NewLabeledSpace(x, y, w, h float64, l Label) *Space {
	rect := NewRect(x, y, w, h)
	return &Space{
		Location: rect,
		Label:    l,
		Type:     NONE,
	}
}

// NewFullSpace returns a space with both a label and a caller id
func NewFullSpace(x, y, w, h float64, l Label, cID event.CID) *Space {
	rect := NewRect(x, y, w, h)
	return &Space{
		rect,
		l,
		cID,
		CID, // todo: This is hard to read as distinct from cID
		// todo: a way to generate non-CID typed spaces that isn't
		// package specific (see render/particle)
	}
}

// NewRect is a wrapper around rtreego.NewRect,
// casting the given x,y to an rtreego.Point.
// Used to not expose rtreego.Point to the user.
func NewRect(x, y, w, h float64) *rtreego.Rect {
	rect, err := rtreego.NewRect(rtreego.Point{x, y}, [3]float64{w, h, 1})
	if err != nil {
		dlog.Error(err)
	}
	return &rect
}
