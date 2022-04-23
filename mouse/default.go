package mouse

import "github.com/oakmound/oak/v4/collision"

// DefaultTree is a collision tree intended to be used by default if no other
// is instantiated. Methods on a collision tree are duplicated as functions
// in this package, so `tree.Add(...)` can instead be `mouse.Add(...)` if
// the codebase is coordinated to just use the default tree.
var (
	DefaultTree = collision.NewTree()
)

// Clear resets the default collision tree
func Clear() {
	DefaultTree.Clear()
}

// Add adds a set of spaces to the rtree
func Add(sps ...*collision.Space) {
	DefaultTree.Add(sps...)
}

// Remove removes a space from the rtree
func Remove(sps ...*collision.Space) {
	DefaultTree.Remove(sps...)
}

// UpdateSpace resets a space's location to a given
// rtreego.Rect.
// This is not an operation on a space because
// a space can exist in multiple rtrees.
func UpdateSpace(x, y, w, h float64, s *collision.Space) error {
	return DefaultTree.UpdateSpace(x, y, w, h, s)
}

// ShiftSpace adds x and y to a space and updates its position
// in the collision rtree that should not be a package global
func ShiftSpace(x, y float64, s *collision.Space) error {
	return DefaultTree.ShiftSpace(x, y, s)
}

// Hits returns the set of spaces which are colliding
// with the passed in space.
func Hits(sp *collision.Space) []*collision.Space {
	return DefaultTree.Hits(sp)
}

// HitLabel acts like hits, but reutrns the first space within hits
// that matches one of the input labels
func HitLabel(sp *collision.Space, labels ...collision.Label) *collision.Space {
	return DefaultTree.HitLabel(sp, labels...)
}
