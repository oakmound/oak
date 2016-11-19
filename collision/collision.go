// Package collision uses an rtree to track rectangles
// and their intersections.
package collision

import (
	"github.com/Sythe2o0/rtreego"
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

func (cp CollisionPoint) IsNil() bool {
	return cp.Zone == nil
}

func Init() {
	rt = rtreego.NewTree(20, 40)
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

func HitLabel(sp *Space, labels ...int) bool {
	results := rt.SearchIntersect(sp.Bounds())
	for _, v := range results {
		for _, label := range labels {
			if v.(*Space) != sp && v.(*Space).Label == label {
				return true
			}
		}
	}
	return false
}
