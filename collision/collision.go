package collision

import (
	//"fmt"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
	"github.com/dhconnelly/rtreego"
	"image/color"
	"log"
)

var (
	rt *rtreego.Rtree
)

type Space struct {
	Location *rtreego.Rect
	cID      event.CID
}

func (s Space) Bounds() *rtreego.Rect {
	return s.Location
}

func Init() {
	rt = rtreego.NewTree(2, 20, 40)
}

func Clear() {
	Init()
}

func Add(sp Space) {
	rt.Insert(sp)
}

func Remove(sp Space) {
	rt.Delete(sp)
}

func UpdateSpace(x, y, w, h float64, s Space) *rtreego.Rect {
	x -= w
	y -= h
	loc := NewRect(x, y, w, h)
	Update(s, loc)
	return loc
}

func Update(s Space, loc *rtreego.Rect) {
	rt.Delete(s)
	s.Location = loc
	rt.Insert(s)
}

func Hits(sp Space) []Space {
	results := rt.SearchIntersect(sp.Bounds())
	out := make([]Space, len(results))
	for index, v := range results {
		out[index] = v.(Space)
	}
	return out
}

func NewSpace(x, y, w, h float64, cID event.CID) Space {
	x -= w
	y -= h
	rect := NewRect(x, y, w, h)
	return Space{
		rect,
		cID,
	}
}

func NewUnassignedSpace(x, y, w, h float64) Space {
	render.DrawColor(color.RGBA{128, 0, 128, 255}, x, y, w, h, 10)
	x -= w
	y -= h
	rect := NewRect(x, y, w, h)
	return Space{Location: rect}
}

func NewRect(x, y, w, h float64) *rtreego.Rect {
	rect, err := rtreego.NewRect(rtreego.Point{x, y}, []float64{w, h})
	if err != nil {
		log.Fatal(err)
	}
	return rect
}
