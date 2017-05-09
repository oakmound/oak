package mouse

import (
	"bitbucket.org/oakmoundstudio/oak/collision"
	"bitbucket.org/oakmoundstudio/oak/event"
	"errors"
)

type CollisionPhase struct {
	OnCollisionS *collision.Space

	wasTouching bool
}

func (cp *CollisionPhase) getCollisionPhase() *CollisionPhase {
	return cp
}

type collisionPhase interface {
	getCollisionPhase() *CollisionPhase
}

func PhaseCollision(s *collision.Space) error {
	switch t := event.GetEntity(int(s.CID)).(type) {
	case collisionPhase:
		oc := t.getCollisionPhase()
		oc.OnCollisionS = s
		s.CID.Bind(phaseCollisionEnter, "EnterFrame")
		return nil
	}
	return errors.New("This space's entity does not implement OnCollisionyThing")
}

func phaseCollisionEnter(id int, nothing interface{}) int {
	e := event.GetEntity(id).(collisionPhase)
	oc := e.getCollisionPhase()

	if oc.OnCollisionS.Contains(LastMouseEvent.ToSpace()) {
		if !oc.wasTouching {
			event.CID(id).Trigger("MouseCollisionStart", LastMouseEvent)
			oc.wasTouching = true
		}
	} else {
		if oc.wasTouching {
			event.CID(id).Trigger("MouseCollisionStop", LastMouseEvent)
			oc.wasTouching = false
		}

	}
	return 0
}
