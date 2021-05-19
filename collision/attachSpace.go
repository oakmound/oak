package collision

import (
	"errors"

	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/physics"
)

// An AttachSpace is a composable struct that provides attachment
// functionality for entities. An entity with AttachSpace can have its
// associated space passed into Attach with the vector the space should
// be attached to.
// Example usage: Any moving character with a collision space. When
// moving the character around by the vector passed in to Attach, the space
// will move with it.
type AttachSpace struct {
	follow     physics.Vector
	aSpace     **Space
	tree       *Tree
	offX, offY float64
}

func (as *AttachSpace) getAttachSpace() *AttachSpace {
	return as
}

type attachSpace interface {
	getAttachSpace() *AttachSpace
}

// Attach attaches v to the given space with optional x,y offsets. See AttachSpace.
func Attach(v physics.Vector, s *Space, tree *Tree, offsets ...float64) error {
	if t, ok := s.CID.E().(attachSpace); ok {
		as := t.getAttachSpace()
		as.aSpace = &s
		as.follow = v
		as.tree = tree
		if as.tree == nil {
			as.tree = DefaultTree
		}
		s.CID.Bind(event.Enter, attachSpaceEnter)
		if len(offsets) > 0 {
			as.offX = offsets[0]
			if len(offsets) > 1 {
				as.offY = offsets[1]
			}
		}
		return nil
	}
	return errors.New("this space's entity is not composed of AttachSpace")
}

// Detach removes the attachSpaceEnter binding from an entity composed with
// AttachSpace
func Detach(s *Space) error {
	en := s.CID.E()
	if _, ok := en.(attachSpace); ok {
		event.UnbindBindable(
			event.UnbindOption{
				Event: event.Event{
					Name:     event.Enter,
					CallerID: s.CID,
				},
				Fn: attachSpaceEnter,
			},
		)
		return nil
	}
	return errors.New("this space's entity is not composed of AttachSpace")
}

func attachSpaceEnter(id event.CID, _ interface{}) int {
	as := id.E().(attachSpace).getAttachSpace()
	x, y := as.follow.X()+as.offX, as.follow.Y()+as.offY
	if x != (*as.aSpace).X() ||
		y != (*as.aSpace).Y() {

		// If this was a nil pointer it would have already crashed but as of release 2.2.0
		// this could error from the space to delete not existing in the rtree.
		as.tree.UpdateSpace(x, y, (*as.aSpace).GetW(), (*as.aSpace).GetH(), *as.aSpace)
	}
	return 0
}
