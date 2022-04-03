package entities

import (
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/physics"
	"github.com/oakmound/oak/v3/render"
)

// Interactive parallels Moving, but for Reactive instead of Solid
type Interactive struct {
	Reactive
	vMoving
}

// NewInteractive returns a new Interactive
func NewInteractive(x, y, w, h float64, r render.Renderable, tree *collision.Tree,
	cid event.CallerID, friction float64) *Interactive {

	i := &Interactive{}
	if cid == 0 {
		i.CallerID = event.DefaultCallerMap.Register(i)
	} else {
		i.CallerID = cid
	}
	i.Reactive = *NewReactive(x, y, w, h, r, tree, i.CallerID)
	i.vMoving = vMoving{
		Delta:    physics.NewVector(0, 0),
		Speed:    physics.NewVector(0, 0),
		Friction: friction,
	}
	return i
}

func (i *Interactive) CID() event.CallerID {
	return i.CallerID
}
