package collision

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
	"github.com/dhconnelly/rtreego"
	"image/color"
)

type Space struct {
	Location *rtreego.Rect
	Label    int
	CID      event.CID
}

func (s Space) Bounds() *rtreego.Rect {
	return s.Location
}

func NewUnassignedSpace(x, y, w, h float64) Space {
	render.DrawColor(color.RGBA{128, 0, 128, 100}, x, y, w, h, 10)
	rect := NewRect(x, y, w, h)
	return Space{Location: rect}
}

func NewSpace(x, y, w, h float64, cID event.CID) Space {
	rect := NewRect(x, y, w, h)
	return Space{
		rect,
		-1,
		cID,
	}
}

func NewLabeledSpace(x, y, w, h float64, l int) Space {
	rect := NewRect(x, y, w, h)
	return Space{
		Location: rect,
		Label:    l,
	}
}

func NewFullSpace(x, y, w, h float64, l int, cID event.CID) Space {
	rect := NewRect(x, y, w, h)
	return Space{
		rect,
		l,
		cID,
	}
}
