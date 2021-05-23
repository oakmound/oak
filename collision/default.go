package collision

// DefaultTree is a collision tree intended to be used by default if no other
// is instantiated. Methods on a collision tree are duplicated as functions
// in this package, so `tree.Add(...)` can instead be `collision.Add(...)` if
// the codebase is coordinated to just use the default tree.
var (
	DefaultTree = NewTree()
)

// Clear resets the default tree's contents
func Clear() {
	DefaultTree.Clear()
}

// Add adds a set of spaces to the rtree
func Add(sps ...*Space) {
	DefaultTree.Add(sps...)
}

// Remove removes a space from the rtree
func Remove(sps ...*Space) {
	DefaultTree.Remove(sps...)
}

// UpdateSpace resets a space's location to a given rect.
func UpdateSpace(x, y, w, h float64, s *Space) error {
	return DefaultTree.UpdateSpace(x, y, w, h, s)
}

// ShiftSpace adds x and y to a space and updates its position
// in the collision rtree that should not be a package global
func ShiftSpace(x, y float64, s *Space) error {
	return DefaultTree.ShiftSpace(x, y, s)
}

// Hits returns the set of spaces which are colliding
// with the passed in space.
func Hits(sp *Space) []*Space {
	return DefaultTree.Hits(sp)
}

// HitLabel acts like hits, but reutrns the first space within hits
// that matches one of the input labels
func HitLabel(sp *Space, labels ...Label) *Space {
	return DefaultTree.HitLabel(sp, labels...)
}

// Update updates this space with the default rtree
func (s *Space) Update(x, y, w, h float64) error {
	return DefaultTree.UpdateSpace(x, y, w, h, s)
}

// SetDim sets the dimensions of the space in the default rtree
func (s *Space) SetDim(w, h float64) error {
	return s.Update(s.X(), s.Y(), w, h)
}

// UpdateLabel changes the label behind this space and resets
// it in the default rtree
func (s *Space) UpdateLabel(classtype Label) {
	DefaultTree.UpdateLabel(classtype, s)
}
