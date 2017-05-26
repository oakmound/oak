package entities

import (
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/physics"
	"bitbucket.org/oakmoundstudio/oak/render"
)

type Interactive struct {
	Reactive
	vMoving
}

func NewInteractive(x, y, w, h float64, r render.Renderable, cid event.CID, friction float64) Interactive {
	return Interactive{
		Reactive: NewReactive(x, y, w, h, r, cid),
		vMoving: vMoving{
			Delta:    physics.NewVector(0, 0),
			Speed:    physics.NewVector(0, 0),
			Friction: friction,
		},
	}
}

func (iv *Interactive) Init() event.CID {
	cID := event.NextID(iv)
	iv.CID = cID
	return cID
}
