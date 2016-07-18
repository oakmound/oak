package collision

import (
	//"fmt"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
	"github.com/dhconnelly/rtreego"
	"image/color"
	"log"
	"math"
	"strconv"
)

var (
	rt *rtreego.Rtree
)

type Space struct {
	Location *rtreego.Rect
	Label    string
	cID      event.CID
}

type CollisionPoint struct {
	Zone *Space
	X    float64
	y    float64
}

func (s Space) Bounds() *rtreego.Rect {
	return s.Location
}

func Init() {
	rt = rtreego.NewTree(2, 20, 40)
	//data structure for raycast
}

func Clear() {
	Init()
}

func Add(sp Space) {
	rt.Insert(sp)
	//space with type
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
		strconv.Itoa(int(cID)),
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

func NewRect(x, y, w, h float64) *rtreego.Rect {
	rect, err := rtreego.NewRect(rtreego.Point{x, y}, []float64{w, h})
	if err != nil {
		log.Fatal(err)
	}
	return rect
}

func RayCast(x, y, degrees, length float64) []CollisionPoint {
	results := []CollisionPoint{}
	resultHash := make(map[Space]bool)

	s := math.Sin(degrees * math.Pi / 180)
	c := math.Cos(degrees * math.Pi / 180)
	for i := 0.0; i < length; i++ {
		loc := NewRect(x, y, .1, .1)

		next := rt.SearchIntersect(loc)

		for k := 0; k < len(next); k++ {
			nx := (next[k].(Space))
			nx_p := &nx
			if _, ok := resultHash[nx]; !ok {
				resultHash[nx] = true
				results = append(results, CollisionPoint{nx_p, x, y})
			}
		}
		x += c
		y += s
	}
	return results
}

func RayCastSingle(x, y, degrees, length float64, invalidIDS []event.CID) CollisionPoint {

	s := math.Sin(degrees * math.Pi / 180)
	c := math.Cos(degrees * math.Pi / 180)
	for i := 0.0; i < length; i++ {
		loc := NewRect(x, y, .1, .1)
		next := rt.SearchIntersect(loc)
	output:
		for k := 0; k < len(next); k++ {
			nx := (next[k].(Space))
			nx_p := &nx
			for e := 0; e < len(invalidIDS); e++ {
				if nx_p.cID == invalidIDS[e] {

					continue output
				}
			}
			return CollisionPoint{nx_p, x, y}
		}
		x += c
		y += s

	}
	return CollisionPoint{}
}
