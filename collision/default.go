package collision

// There's a default collision tree you can access via collision.func
// as opposed to tree.func.
var (
	DefTree *Tree
)

func init() {
	// This won't error so long as DefMinChildren < DefMaxChildren
	DefTree, _ = NewTree()
}

// Clear resets the default tree's contents
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

// Update updates this space with the default rtree
func (s *Space) Update(x, y, w, h float64) error {
	return DefTree.UpdateSpace(x, y, w, h, s)
}

// SetDim sets the dimensions of the space in the default rtree
func (s *Space) SetDim(w, h float64) error {
	return s.Update(s.X(), s.Y(), w, h)
}

// UpdateLabel changes the label behind this space and resets
// it in the default rtree
func (s *Space) UpdateLabel(classtype Label) {
	DefTree.Remove(s)
	s.Label = classtype
	DefTree.Add(s)
}
