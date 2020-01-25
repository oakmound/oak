package entities

import (
	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/physics"
	"github.com/oakmound/oak/v2/render"
)

// Interactive parallels Moving, but for Reactive instead of Solid
type Interactive struct {
	Reactive
	vMoving
}

// NewInteractive returns a new Interactive
func NewInteractive(x, y, w, h float64, r render.Renderable, tree *collision.Tree,
	cid event.CID, friction float64) *Interactive {

	i := Interactive{}
	cid = cid.Parse(&i)
	i.Reactive = *NewReactive(x, y, w, h, r, tree, cid)
	i.vMoving = vMoving{
		Delta:    physics.NewVector(0, 0),
		Speed:    physics.NewVector(0, 0),
		Friction: friction,
	}
	return &i
}

// Init satisfies event.Entity
func (iv *Interactive) Init() event.CID {
	cID := event.NextID(iv)
	iv.CID = cID
	return cID
}
