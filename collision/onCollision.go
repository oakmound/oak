package collision

import (
	"errors"

	"github.com/oakmound/oak/v3/event"
)

// A Phase is a struct that other structs who want to use PhaseCollision
// should be composed of
type Phase struct {
	OnCollisionS *Space
	tree         *Tree
	bus          event.Handler
	// If allocating maps becomes an issue
	// we can have two constant maps that we
	// switch between on alternating frames
	Touching map[Label]bool
}

func (cp *Phase) getCollisionPhase() *Phase {
	return cp
}

type collisionPhase interface {
	getCollisionPhase() *Phase
}

// PhaseCollision binds to the entity behind the space's CID so that it will
// receive CollisionStart and CollisionStop events, appropriately when
// entities begin to collide or stop colliding with the space.
// If tree is nil, it uses DefTree
func PhaseCollision(s *Space, tree *Tree) error {
	return PhaseCollisionWithBus(s, tree, event.DefaultBus, event.DefaultCallerMap)
}

// PhaseCollisionWithBus allows for a non-default bus and non-default entity mapping
// in a phase collision binding.
func PhaseCollisionWithBus(s *Space, tree *Tree, bus event.Handler, entities *event.CallerMap) error {
	en := entities.GetEntity(s.CID)
	if cp, ok := en.(collisionPhase); ok {
		oc := cp.getCollisionPhase()
		oc.OnCollisionS = s
		oc.tree = tree
		oc.bus = bus
		if oc.tree == nil {
			oc.tree = DefaultTree
		}
		bus.UnsafeBind(event.Enter.UnsafeEventID, s.CID, phaseCollisionEnter)
		return nil
	}
	return errors.New("This space's entity does not implement collisionPhase")
}

// CollisionStart/Stop: when a PhaseCollision entity starts/stops touching some label.
var (
	Start = event.RegisterEvent[Label]()
	Stop  = event.RegisterEvent[Label]()
)

func phaseCollisionEnter(id event.CallerID, handler event.Handler, _ interface{}) event.Response {
	e := handler.GetCallerMap().GetEntity(id).(collisionPhase)
	oc := e.getCollisionPhase()

	// check hits
	hits := oc.tree.Hits(oc.OnCollisionS)
	newTouching := map[Label]bool{}

	// if any are new, trigger on collision
	for _, h := range hits {
		l := h.Label
		if _, ok := oc.Touching[l]; !ok {
			event.TriggerForCallerOn(oc.bus, id, Start, l)
		}
		newTouching[l] = true
	}

	// if we lost any, trigger off collision
	for l := range oc.Touching {
		if _, ok := newTouching[l]; !ok {
			event.TriggerForCallerOn(handler, id, Stop, l)
		}
	}

	oc.Touching = newTouching

	return 0
}
