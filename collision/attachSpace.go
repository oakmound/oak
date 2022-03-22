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
	binding    event.Binding
}

func (as *AttachSpace) getAttachSpace() *AttachSpace {
	return as
}

func (as *AttachSpace) CID() event.CallerID {
	return (*as.aSpace).CID
}

var _ attachSpace = &AttachSpace{}

type attachSpace interface {
	event.Caller
	getAttachSpace() *AttachSpace
}

// Attach attaches v to the given space with optional x,y offsets. See AttachSpace.
func Attach(v physics.Vector, s *Space, tree *Tree, offsets ...float64) error {
	en := event.DefaultCallerMap.GetEntity(s.CID)
	if t, ok := en.(attachSpace); ok {
		as := t.getAttachSpace()
		as.aSpace = &s
		as.follow = v
		as.tree = tree
		if as.tree == nil {
			as.tree = DefaultTree
		}
		as.binding = event.Bind(event.DefaultBus, event.Enter, t, attachSpaceEnter)
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
	en := event.DefaultCallerMap.GetEntity(s.CID)
	if as, ok := en.(attachSpace); ok {
		as.getAttachSpace().binding.Unbind()
		return nil
	}
	return errors.New("this space's entity is not composed of AttachSpace")
}

func attachSpaceEnter(asIface attachSpace, _ event.EnterPayload) event.Response {
	as := asIface.(attachSpace).getAttachSpace()
	x, y := as.follow.X()+as.offX, as.follow.Y()+as.offY
	if x != (*as.aSpace).X() ||
		y != (*as.aSpace).Y() {
		as.tree.UpdateSpace(x, y, (*as.aSpace).GetW(), (*as.aSpace).GetH(), *as.aSpace)
	}
	return 0
}
