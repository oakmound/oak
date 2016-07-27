package collision

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
	"github.com/dhconnelly/rtreego"
	"image/color"
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

func NewUnassignedSpace(x, y, w, h float64) *Space {
	render.DrawColor(color.RGBA{128, 0, 128, 100}, x, y, w, h, 10)
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
