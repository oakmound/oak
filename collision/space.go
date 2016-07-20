package collision

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
	"github.com/dhconnelly/rtreego"
	"image/color"
	"strconv"
)

type Space struct {
	Location *rtreego.Rect
	Label    string
	cID      event.CID
}

func (s Space) Bounds() *rtreego.Rect {
	return s.Location
}

func NewUnassignedSpace(x, y, w, h float64) Space {
	render.DrawColor(color.RGBA{128, 0, 128, 100}, x, y, w, h, 10)
	x -= w
	y -= h
	rect := NewRect(x, y, w, h)
	return Space{Location: rect}
}

func NewSpace(x, y, w, h float64, cID event.CID) Space {
	x -= w
	y -= h
	rect := NewRect(x, y, w, h)
	return Space{
		rect,
		strconv.Itoa(int(cID)),
		cID,
	}
}

func NewLabeledSpace(x, y, w, h float64, s string) Space {
	x -= w
	y -= h
	rect := NewRect(x, y, w, h)
	return Space{
		Location: rect,
		Label:    s,
	}
}

func NewFullSpace(x, y, w, h float64, s string, cID event.CID) Space {
	x -= w
	y -= h
	rect := NewRect(x, y, w, h)
	return Space{
		rect,
		s,
		cID,
	}
}
