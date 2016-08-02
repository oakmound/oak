// Package collision uses an rtree to track rectangles
// and their intersections.
package collision

import (
	"github.com/dhconnelly/rtreego"
	"log"
)

var (
	rt *rtreego.Rtree
)

// A CollisionPoint is a specific point where
// collision occured and a zone to identify
// what was collided with.
type CollisionPoint struct {
	Zone *Space
	X, Y float64
}

func Init() {
	rt = rtreego.NewTree(2, 20, 40)
}

func Clear() {
	Init()
}

func Add(sp *Space) {
	rt.Insert(sp)
}

func Remove(sp *Space) {
	rt.Delete(sp)
}

// Update resets a space's location to a given
// rtreego.Rect.
// This is not an operation on a space because
// a space can exist in multiple rtrees.
func UpdateSpace(x, y, w, h float64, s *Space) {
	loc := NewRect(x, y, w, h)
	rt.Delete(s)
	s.Location = loc
	rt.Insert(s)
}

// Hits returns the set of spaces which are colliding
// with the passed in space.
func Hits(sp *Space) []*Space {
	results := rt.SearchIntersect(sp.Bounds())
	out := make([]*Space, len(results))
	for index, v := range results {
		out[index] = v.(*Space)
	}
	return out
}

// NewRect is a wrapper around rtreego.NewRect,
// casting the given x,y to an rtreego.Point.
// Used to not expose rtreego.Point to the user.
func NewRect(x, y, w, h float64) *rtreego.Rect {
	rect, err := rtreego.NewRect(rtreego.Point{x, y}, []float64{w, h})
	if err != nil {
		log.Fatal(err)
	}
	return rect
}
