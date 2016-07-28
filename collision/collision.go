// Package collision uses an rtree to track rectangles
// and their intersections.
package collision

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"github.com/dhconnelly/rtreego"
	"log"
	"math"
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

// RayCast returns the set of points where a line
// from x,y going at a certain angle, for a certain length, intersects
// with existing rectangles in the rtree.
// It converts the ray into a series of points which are themselves
// used to check collision at a miniscule width and height.
func RayCast(x, y, degrees, length float64) []CollisionPoint {
	results := []CollisionPoint{}
	resultHash := make(map[*Space]bool)

	s := math.Sin(degrees * math.Pi / 180)
	c := math.Cos(degrees * math.Pi / 180)
	for i := 0.0; i < length; i++ {
		loc := NewRect(x, y, .1, .1)

		next := rt.SearchIntersect(loc)

		for k := 0; k < len(next); k++ {
			nx := (next[k].(*Space))
			if _, ok := resultHash[nx]; !ok {
				resultHash[nx] = true
				results = append(results, CollisionPoint{nx, x, y})
			}
		}
		x += c
		y += s
	}
	return results
}

// RatCastSingle acts as RayCast, but it returns only the first collision
// that the generated ray intersects, ignoring entities
// in the given invalidIDs list.
// Example Use case: shooting a bullet, hitting the first thing that isn't yourself.
func RayCastSingle(x, y, degrees, length float64, invalidIDS []event.CID) CollisionPoint {

	s := math.Sin(degrees * math.Pi / 180)
	c := math.Cos(degrees * math.Pi / 180)
	for i := 0.0; i < length; i++ {
		loc := NewRect(x, y, .1, .1)
		next := rt.SearchIntersect(loc)
	output:
		for k := 0; k < len(next); k++ {
			nx := (next[k].(*Space))
			for e := 0; e < len(invalidIDS); e++ {
				if nx.CID == invalidIDS[e] {
					continue output
				}
			}
			return CollisionPoint{nx, x, y}
		}
		x += c
		y += s

	}
	return CollisionPoint{}
}
