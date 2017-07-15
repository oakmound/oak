package collision

import (
	"errors"
	"sync"

	"github.com/Sythe2o0/rtreego"
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
	*rtreego.Rtree
	sync.Mutex
	minChildren, maxChildren int
}

// NewTree returns a new collision Tree
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
		Rtree:       rtreego.NewTree(minChildren, maxChildren),
		minChildren: minChildren,
		maxChildren: maxChildren,
		Mutex:       sync.Mutex{},
	}, nil
}

// Clear resets a tree's contents to be empty
func (t *Tree) Clear() {
	t.Rtree = rtreego.NewTree(t.minChildren, t.maxChildren)
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

// Remove removes spaces from the rtree
// returns the number of spaces removed
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

// UpdateSpace resets a space's location to a given
// rtreego.Rect.
// This is not an operation on a space because
// a space can exist in multiple rtrees.
func (t *Tree) UpdateSpace(x, y, w, h float64, s *Space) error {
	if s == nil {
		return errors.New("Input space was nil")
	}
	loc := NewRect(x, y, w, h)
	t.Lock()
	t.Delete(s)
	s.Location = loc
	t.Insert(s)
	t.Unlock()
	return nil
}

// ShiftSpace adds x and y to a space and updates its position
func (t *Tree) ShiftSpace(x, y float64, s *Space) error {
	x = x + s.GetX()
	y = y + s.GetY()
	return t.UpdateSpace(x, y, s.GetW(), s.GetH(), s)
}

// Hits returns the set of spaces which are colliding
// with the passed in space.
func (t *Tree) Hits(sp *Space) []*Space {
	// Eventually we'll expose SearchIntersect for use cases where you
	// want to see if you intersect yourself
	results := t.SearchIntersect(sp.Bounds())
	out := make([]*Space, len(results))
	hitSelf := -1
	for i, v := range results {
		if v.(*Space) == sp {
			hitSelf = i
		}
		out[i] = v.(*Space)
	}
	if hitSelf != -1 {
		out[hitSelf], out[len(out)-1] = out[len(out)-1], out[hitSelf]
		return out[:len(out)-1]
	}
	return out
}

// HitLabel acts like hits, but returns the first space within hits
// that matches one of the input labels
func (t *Tree) HitLabel(sp *Space, labels ...Label) *Space {
	results := t.SearchIntersect(sp.Bounds())
	for _, v := range results {
		for _, label := range labels {
			if v.(*Space) != sp && v.(*Space).Label == label {
				return v.(*Space)
			}
		}
	}
	return nil
}

// Hit is an experimental new syntax that probably has performance hits
// relative to Hits/HitLabel, see filters.go
func (t *Tree) Hit(sp *Space, fs ...Filter) []*Space {
	iresults := t.SearchIntersect(sp.Bounds())
	results := make([]*Space, len(iresults))
	for i, v := range iresults {
		results[i] = v.(*Space)
	}
	for _, f := range fs {
		if len(results) == 0 {
			return results
		}
		results = f(results)
	}
	return results
}
