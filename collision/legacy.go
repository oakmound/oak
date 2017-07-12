package collision

import "log"

// There's a default collision tree you can access via collision.func
// as opposed to tree.func
var (
	rt *Tree
)

func init() {
	var err error
	rt, err = NewTree()
	if err != nil {
		log.Fatal(err)
	}
}

// Clear just calls init.
func Clear() {
	rt.Clear()
}

// Add adds a set of spaces to the rtree
func Add(sps ...*Space) {
	rt.Add(sps...)
}

// Remove removes a space from the rtree
func Remove(sps ...*Space) {
	rt.Remove(sps...)
}

// UpdateSpace resets a space's location to a given
// rtreego.Rect.
// This is not an operation on a space because
// a space can exist in multiple rtrees.
func UpdateSpace(x, y, w, h float64, s *Space) {
	rt.UpdateSpace(x, y, w, h, s)
}

// ShiftSpace adds x and y to a space and updates its position
// in the collision rtree that should not be a package global
func ShiftSpace(x, y float64, s *Space) {
	rt.ShiftSpace(x, y, s)
}

// Hits returns the set of spaces which are colliding
// with the passed in space.
func Hits(sp *Space) []*Space {
	return rt.Hits(sp)
}

// HitLabel acts like hits, but reutrns the first space within hits
// that matches one of the input labels
func HitLabel(sp *Space, labels ...Label) *Space {
	return rt.HitLabel(sp, labels...)
}

// Update updates this space with the legacy rtree
func (s *Space) Update(x, y, w, h float64) {
	rt.UpdateSpace(x, y, w, h, s)
}

// UpdateLabel changes the label behind this space and resets
// it in the legacy rtree
func (s *Space) UpdateLabel(classtype Label) {
	rt.Remove(s)
	s.Label = classtype
	rt.Add(s)
}
