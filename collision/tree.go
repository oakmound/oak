package collision

import (
	"errors"
	"sync"

	"github.com/oakmound/oak/v2/alg/floatgeom"
	"github.com/oakmound/oak/v2/oakerr"
)

var (
	// DefaultMaxChildren is the maximum number of children allowed
	// on a node in the collision tree when NewTree is called without
	// a maximum number of children.
	DefaultMaxChildren = 40
	// DefaultMinChildren is the minimum number of children allowed
	// on a node in the collision tree when NewTree is called without
	// a minimum number of children.
	DefaultMinChildren = 20
)

// A Tree provides a space for managing collisions between rectangles
type Tree struct {
	*Rtree
	sync.Mutex
}

// NewTree returns a new collision Tree. The first argument will be used
// as the minimum children per tree node. The second will be the maximum
// children per tree node. Further arguments are ignored. If less than two
// arguments are given, DefaultMinChildren and DefaultMaxChildren will be
// used.
func NewTree(children ...int) (*Tree, error) {
	minChildren := DefaultMinChildren
	maxChildren := DefaultMaxChildren
	if len(children) > 0 {
		minChildren = children[0]
		if len(children) > 1 {
			maxChildren = children[1]
		}
	}
	if minChildren > maxChildren {
		return nil, errors.New("MaxChildren must exceed MinChildren")
	}
	return &Tree{
		Rtree: newTree(minChildren, maxChildren),
		Mutex: sync.Mutex{},
	}, nil
}

// Clear resets a tree's contents to be empty
func (t *Tree) Clear() {
	t.Rtree = newTree(t.Rtree.MinChildren, t.Rtree.MaxChildren)
}

// Add adds a set of spaces to the rtree
func (t *Tree) Add(sps ...*Space) {
	t.Lock()
	for _, sp := range sps {
		if sp != nil {
			t.Insert(sp)
		}
	}
	t.Unlock()
}

// Remove removes spaces from the rtree and
// returns the number of spaces removed.
func (t *Tree) Remove(sps ...*Space) int {
	removed := 0
	t.Lock()
	for _, sp := range sps {
		if sp != nil {
			if t.Delete(sp) {
				removed++
			}
		}
	}
	t.Unlock()
	return removed
}

// UpdateLabel will set the input space's label. DEPRECATED. Just set
// the Label field on the Space pointer.
func (t *Tree) UpdateLabel(classtype Label, s *Space) {
	s.Label = classtype
}

// ErrNotExist is returned by methods on spaces
// when the space to update or act on did not exist
var ErrNotExist = errors.New("Space did not exist to update")

// UpdateSpace resets a space's location to a given
// rtreego.Rect.
// This is not an operation on a space because
// a space can exist in multiple rtrees.
func (t *Tree) UpdateSpace(x, y, w, h float64, s *Space) error {
	loc := NewRect(x, y, w, h)
	return t.UpdateSpaceRect(loc, s)
}

// UpdateSpaceRect acts as UpdateSpace, but takes in a rectangle instead
// of four distinct arguments.
func (t *Tree) UpdateSpaceRect(rect floatgeom.Rect3, s *Space) error {
	if s == nil {
		return oakerr.NilInput{InputName: "s"}
	}
	t.Lock()
	deleted := t.Delete(s)
	if !deleted {
		t.Unlock()
		return ErrNotExist
	}
	s.Location = rect
	t.Insert(s)
	t.Unlock()
	return nil
}

// ShiftSpace adds x and y to a space and updates its position
func (t *Tree) ShiftSpace(x, y float64, s *Space) error {
	if s == nil {
		return oakerr.NilInput{InputName: "s"}
	}
	x += s.X()
	y += s.Y()
	return t.UpdateSpace(x, y, s.GetW(), s.GetH(), s)
}

// Hits returns the set of spaces which are colliding
// with the passed in space. All spaces collide with
// themselves, if they exist in the tree, but self-collision
// will not be reported by Hits.
func (t *Tree) Hits(sp *Space) []*Space {
	results := t.SearchIntersect(sp.Bounds())
	hitSelf := -1
	i := 0
	for i < len(results) {
		// Todo: replicate getting nils again, its not happening anymore
		// without unexported field modification
		if results[i] == nil {
			results = append(results[:i], results[i+1:]...)
		} else {
			i++
		}
	}
	out := make([]*Space, len(results))
	for i, v := range results {
		if v == sp {
			hitSelf = i
		}
		out[i] = v
	}
	if hitSelf != -1 {
		out[hitSelf], out[len(out)-1] = out[len(out)-1], out[hitSelf]
		return out[:len(out)-1]
	}
	return out
}

// HitLabel acts like Hits, but returns the first space within hits
// that matches one of the input labels. HitLabel can return the same
// space that is passed into it, if that space has a label in the set of
// accepted labels.
func (t *Tree) HitLabel(sp *Space, labels ...Label) *Space {
	results := t.SearchIntersect(sp.Bounds())
	for _, v := range results {
		for _, label := range labels {
			if v != sp && v.Label == label {
				return v
			}
		}
	}
	return nil
}

// Hit is an experimental new syntax that probably has performance hits
// relative to Hits/HitLabel, see filters.go
func (t *Tree) Hit(sp *Space, fs ...Filter) []*Space {
	results := t.SearchIntersect(sp.Bounds())
	for _, f := range fs {
		if len(results) == 0 {
			return results
		}
		results = f(results)
	}
	return results
}
