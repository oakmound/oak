package mouse

import (
	"log"

	"github.com/oakmound/oak/collision"
)

// There's a default collision tree you can access via collision.func
// as opposed to tree.func
var (
	DefTree *collision.Tree
)

func init() {
	var err error
	DefTree, err = collision.NewTree()
	if err != nil {
		log.Fatal(err)
	}
}

// Clear just calls init.
func Clear() {
	DefTree.Clear()
}

// Add adds a set of spaces to the rtree
func Add(sps ...*collision.Space) {
	DefTree.Add(sps...)
}

// Remove removes a space from the rtree
func Remove(sps ...*collision.Space) {
	DefTree.Remove(sps...)
}

// UpdateSpace resets a space's location to a given
// rtreego.Rect.
// This is not an operation on a space because
// a space can exist in multiple rtrees.
func UpdateSpace(x, y, w, h float64, s *collision.Space) error {
	return DefTree.UpdateSpace(x, y, w, h, s)
}

// ShiftSpace adds x and y to a space and updates its position
// in the collision rtree that should not be a package global
func ShiftSpace(x, y float64, s *collision.Space) error {
	return DefTree.ShiftSpace(x, y, s)
}

// Hits returns the set of spaces which are colliding
// with the passed in space.
func Hits(sp *collision.Space) []*collision.Space {
	return DefTree.Hits(sp)
}

// HitLabel acts like hits, but reutrns the first space within hits
// that matches one of the input labels
func HitLabel(sp *collision.Space, labels ...collision.Label) *collision.Space {
	return DefTree.HitLabel(sp, labels...)
}
