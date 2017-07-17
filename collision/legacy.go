package collision

import (
	"log"

	"github.com/oakmound/oak/event"
)

// There's a default collision tree you can access via collision.func
// as opposed to tree.func. This is considered a legacy set of features,
// because the benefit to the API is minimal in exchange for a much harder
// to use collision tree. It does make small applications a little shorter.
var (
	DefTree *Tree
)

func init() {
	var err error
	DefTree, err = NewTree()
	if err != nil {
		log.Fatal(err)
	}
}

// Clear just calls init.
func Clear() {
	DefTree.Clear()
}

// Add adds a set of spaces to the rtree
func Add(sps ...*Space) {
	DefTree.Add(sps...)
}

// Remove removes a space from the rtree
func Remove(sps ...*Space) {
	DefTree.Remove(sps...)
}

// UpdateSpace resets a space's location to a given
// rtreego.Rect.
// This is not an operation on a space because
// a space can exist in multiple rtrees.
func UpdateSpace(x, y, w, h float64, s *Space) error {
	return DefTree.UpdateSpace(x, y, w, h, s)
}

// ShiftSpace adds x and y to a space and updates its position
// in the collision rtree that should not be a package global
func ShiftSpace(x, y float64, s *Space) error {
	return DefTree.ShiftSpace(x, y, s)
}

// Hits returns the set of spaces which are colliding
// with the passed in space.
func Hits(sp *Space) []*Space {
	return DefTree.Hits(sp)
}

// HitLabel acts like hits, but reutrns the first space within hits
// that matches one of the input labels
func HitLabel(sp *Space, labels ...Label) *Space {
	return DefTree.HitLabel(sp, labels...)
}

// Update updates this space with the legacy rtree
func (s *Space) Update(x, y, w, h float64) error {
	return DefTree.UpdateSpace(x, y, w, h, s)
}

// UpdateLabel changes the label behind this space and resets
// it in the legacy rtree
func (s *Space) UpdateLabel(classtype Label) {
	DefTree.Remove(s)
	s.Label = classtype
	DefTree.Add(s)
}

// RayCast returns the set of points where a line
// from x,y going at a certain angle, for a certain length, intersects
// with existing rectangles in the rtree.
// It converts the ray into a series of points which are themselves
// used to check collision at a miniscule width and height.
func RayCast(x, y, degrees, length float64) []Point {
	return DefTree.RayCast(x, y, degrees, length)
}

// RayCastSingle acts as RayCast, but it returns only the first collision
// that the generated ray intersects, ignoring entities
// in the given invalidIDs list.
// Example Use case: shooting a bullet, hitting the first thing that isn't yourself.
func RayCastSingle(x, y, degrees, length float64, invalidIDS []event.CID) Point {
	return DefTree.RayCastSingle(x, y, degrees, length, invalidIDS)
}

// RayCastSingleLabels acts like RayCastSingle, but only returns elements
// that match one of the input labels
func RayCastSingleLabels(x, y, degrees, length float64, labels ...Label) Point {
	return DefTree.RayCastSingleLabels(x, y, degrees, length, labels...)
}

// RayCastSingleIgnoreLabels is the opposite of Labels, in that it will return
// the first collision point that is not contained in the set of ignore labels
func RayCastSingleIgnoreLabels(x, y, degrees, length float64, labels ...Label) Point {
	return DefTree.RayCastSingleIgnoreLabels(x, y, degrees, length, labels...)
}

// RayCastSingleIgnore is just like ignore labels but also ignores certain
// caller ids
func RayCastSingleIgnore(x, y, degrees, length float64, invalidIDS []event.CID, labels ...Label) Point {
	return DefTree.RayCastSingleIgnore(x, y, degrees, length, invalidIDS, labels...)
}

// ConeCast repeatedly calls RayCast in a cone shape
// ConeCast advances COUNTER-CLOCKWISE
func ConeCast(x, y, angle, angleWidth, rays, length float64) (points []Point) {
	return DefTree.ConeCast(x, y, angle, angleWidth, rays, length)
}

// ConeCastSingle repeatedly calls RayCastSignle in a cone shape
func ConeCastSingle(x, y, angle, angleWidth, rays, length float64, invalidIDS []event.CID) (points []Point) {
	return DefTree.ConeCastSingle(x, y, angle, angleWidth, rays, length, invalidIDS)
}

// ConeCastSingleLabels repeatedly calls RayCastSingleLabels in a cone shape
func ConeCastSingleLabels(x, y, angle, angleWidth, rays, length float64, labels ...Label) (points []Point) {
	return DefTree.ConeCastSingleLabels(x, y, angle, angleWidth, rays, length, labels...)
}
