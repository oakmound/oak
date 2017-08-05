package collision

import (
	"errors"

	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
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
	offX, offY float64
}

func (as *AttachSpace) getAttachSpace() *AttachSpace {
	return as
}

type attachSpace interface {
	getAttachSpace() *AttachSpace
}

// Attach attaches v to the given space with optional x,y offsets. See AttachSpace.
// Attach binds attachSpaceEnter at priority -1. This means that attachSpaceEnter,
// which updates the collision space for an AttachSpace composed entity, will be called
// after all EnterFrame bindings that are bound with .Bind(), but before those that
// are called with .BindPriority(... -2).
func Attach(v physics.Vector, s *Space, offsets ...float64) error {
	if t, ok := event.GetEntity(int(s.CID)).(attachSpace); ok {
		as := t.getAttachSpace()
		as.aSpace = &s
		as.follow = v
		s.CID.BindPriority(attachSpaceEnter, "EnterFrame", -1)
		if len(offsets) > 0 {
			as.offX = offsets[0]
			if len(offsets) > 1 {
				as.offY = offsets[1]
			}
		}
		return nil
	}
	return errors.New("This space's entity is not composed of AttachSpace")
}

// Detach removes the attachSpaceEnter binding from an entity composed with
// AttachSpace
func Detach(s *Space) error {
	switch event.GetEntity(int(s.CID)).(type) {
	case attachSpace:
		// Todo: this syntax is ugly
		// Note UnbindBindable is not a recommended way to unbind things,
		// but is okay here because we know we are not unbinding a closure.
		event.UnbindBindable(
			event.UnbindOption{
				BindingOption: event.BindingOption{
					Event: event.Event{
						Name:     "EnterFrame",
						CallerID: int(s.CID),
					},
					Priority: -1,
				},
				Fn: attachSpaceEnter,
			},
		)
		return nil
	}
	return errors.New("This space's entity is not composed of AttachSpace")
}

// attachSpaceEnter currently uses the default tree, always. Todo: change this,
// see what onCollision does
func attachSpaceEnter(id int, nothing interface{}) int {
	as := event.GetEntity(id).(attachSpace).getAttachSpace()
	x, y := as.follow.X()+as.offX, as.follow.Y()+as.offY
	if x != (*as.aSpace).GetX() ||
		y != (*as.aSpace).GetY() {
		UpdateSpace(x, y, (*as.aSpace).GetW(), (*as.aSpace).GetH(), *as.aSpace)
	}
	return 0
}
