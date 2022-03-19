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
	CallerMap    *event.CallerMap
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
func PhaseCollision(s *collision.Space, callerMap *event.CallerMap, handler event.Handler) error {
	if callerMap == nil {
		callerMap = event.DefaultCallerMap
	}
	en := callerMap.GetEntity(s.CID)
	if cp, ok := en.(collisionPhase); ok {
		oc := cp.getCollisionPhase()
		oc.OnCollisionS = s
		oc.CallerMap = callerMap
		event.Bind(handler, event.Enter, s.CID, phaseCollisionEnter(callerMap, handler))
		return nil
	}
	return errors.New("This space's entity does not implement collisionPhase")
}

// MouseCollisionStart/Stop: see collision Start/Stop, for mouse collision
var (
	Start = event.RegisterEvent[*Event]()
	Stop  = event.RegisterEvent[*Event]()
)

func phaseCollisionEnter(callerMap *event.CallerMap, handler event.Handler) func(id event.CallerID, payload event.EnterPayload) event.Response {
	return func(id event.CallerID, payload event.EnterPayload) event.Response {

		e := callerMap.GetEntity(id).(collisionPhase)
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
}
