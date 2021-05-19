package mouse

import (
	"errors"

	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/event"
)

// CollisionPhase is a component that can be placed into another struct to
// enable PhaseCollision on the struct. See PhaseCollision.
type CollisionPhase struct {
	OnCollisionS *collision.Space
	LastEvent    *Event

	wasTouching bool
}

func (cp *CollisionPhase) getCollisionPhase() *CollisionPhase {
	return cp
}

type collisionPhase interface {
	getCollisionPhase() *CollisionPhase
}

// PhaseCollision binds to the entity behind the space's CID so that it will
// receive MouseCollisionStart and MouseCollisionStop events, appropriately when
// the mouse begins to hover or stops hovering over the input space.
func PhaseCollision(s *collision.Space) error {
	en := s.CID.E()
	if cp, ok := en.(collisionPhase); ok {
		oc := cp.getCollisionPhase()
		oc.OnCollisionS = s
		s.CID.Bind(event.Enter, phaseCollisionEnter)
		return nil
	}
	return errors.New("This space's entity does not implement collisionPhase")
}

// MouseCollisionStart/Stop: see collision Start/Stop, for mouse collision
// Payload: (mouse.Event)
const (
	Start = "MouseCollisionStart"
	Stop  = "MouseCollisionStop"
)

func phaseCollisionEnter(id event.CID, nothing interface{}) int {
	e := id.E().(collisionPhase)
	oc := e.getCollisionPhase()

	// TODO: think about how this can more cleanly work with multiple controllers
	ev := oc.LastEvent
	if ev == nil {
		ev = &LastEvent
	}

	if oc.OnCollisionS.Contains(ev.ToSpace()) {
		if !oc.wasTouching {
			id.Trigger(Start, *ev)
			oc.wasTouching = true
		}
	} else {
		if oc.wasTouching {
			id.Trigger(Stop, *ev)
			oc.wasTouching = false
		}
	}
	return 0
}
