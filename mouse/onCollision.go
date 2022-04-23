package mouse

import (
	"errors"

	"github.com/oakmound/oak/v4/collision"
	"github.com/oakmound/oak/v4/event"
)

// CollisionPhase is a component that can be placed into another struct to
// enable PhaseCollision on the struct. See PhaseCollision.
type CollisionPhase struct {
	OnCollisionS *collision.Space
	CallerMap    *event.CallerMap
	LastEvent    *Event

	wasTouching bool
}

func (cp *CollisionPhase) getCollisionPhase() *CollisionPhase {
	return cp
}

func (cp *CollisionPhase) CID() event.CallerID {
	return cp.OnCollisionS.CID
}

type collisionPhase interface {
	getCollisionPhase() *CollisionPhase
}

// PhaseCollision binds to the entity behind the space's CID so that it will
// receive MouseCollisionStart and MouseCollisionStop events, appropriately when
// the mouse begins to hover or stops hovering over the input space.
func PhaseCollision(s *collision.Space, handler event.Handler) error {
	en := handler.GetCallerMap().GetEntity(s.CID)
	if cp, ok := en.(collisionPhase); ok {
		oc := cp.getCollisionPhase()
		oc.OnCollisionS = s
		oc.CallerMap = handler.GetCallerMap()
		handler.UnsafeBind(event.Enter.UnsafeEventID, s.CID, phaseCollisionEnter)
		return nil
	}
	return errors.New("This space's entity does not implement collisionPhase")
}

// MouseCollisionStart/Stop: see collision Start/Stop, for mouse collision
var (
	Start = event.RegisterEvent[*Event]()
	Stop  = event.RegisterEvent[*Event]()
)

func phaseCollisionEnter(id event.CallerID, handler event.Handler, _ interface{}) event.Response {
	e, ok := handler.GetCallerMap().GetEntity(id).(collisionPhase)
	if !ok {
		return event.ResponseUnbindThisBinding
	}
	oc := e.getCollisionPhase()
	if oc == nil || oc.OnCollisionS == nil {
		return 0
	}

	// TODO: think about how this can more cleanly work with multiple windows
	ev := oc.LastEvent
	if ev == nil {
		ev = &LastEvent
	}
	if ev.StopPropagation {
		return 0
	}

	if oc.OnCollisionS.Contains(ev.ToSpace()) {
		if !oc.wasTouching {
			event.TriggerForCallerOn(handler, id, Start, ev)
			oc.wasTouching = true
		}
	} else {
		if oc.wasTouching {
			event.TriggerForCallerOn(handler, id, Stop, ev)
			oc.wasTouching = false
		}
	}
	return 0
}
