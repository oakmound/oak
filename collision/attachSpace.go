package collision

import (
	"errors"

	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/physics"
)

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

// Attach binds attachSpaceEnter at priority -1
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

func Detach(s *Space) error {
	switch event.GetEntity(int(s.CID)).(type) {
	case attachSpace:
		// Todo: this syntax is terrible
		event.UnbindBindable(
			event.UnbindOption{
				event.BindingOption{
					event.Event{
						"EnterFrame",
						int(s.CID),
					},
					0,
				},
				attachSpaceEnter,
			},
		)
		return nil
	}
	return errors.New("This space's entity is not composed of AttachSpace")
}

func attachSpaceEnter(id int, nothing interface{}) int {
	as := event.GetEntity(id).(attachSpace).getAttachSpace()
	if as.follow.X()+as.offX != (*as.aSpace).GetX() ||
		as.follow.Y()+as.offY != (*as.aSpace).GetY() {
		UpdateSpace(as.follow.X()+as.offX, as.follow.Y()+as.offY, (*as.aSpace).GetW(), (*as.aSpace).GetH(), *as.aSpace)
	}
	return 0
}
