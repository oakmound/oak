package entities

import (
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/physics"
	"bitbucket.org/oakmoundstudio/oak/render"
)

type Interactive struct {
	Reactive
	moving
}

func NewInteractive(x, y, w, h float64, r render.Renderable, cid event.CID) Interactive {
	return Interactive{
		Reactive: NewReactive(x, y, w, h, r, cid),
		moving:   moving{},
	}
}

func (i *Interactive) Init() event.CID {
	cID := event.NextID(i)
	i.CID = cID
	return cID
}

func (i *Interactive) String() string {
	st := "Interactive: \n{"
	st += i.Reactive.String()
	st += "}\n " + i.moving.String()
	return st
}

type InteractVector struct {
	Reactive
	vMoving
}

func NewInteractVector(x, y, w, h float64, r render.Renderable, cid event.CID, friction float64) InteractVector {
	return InteractVector{
		Reactive: NewReactive(x, y, w, h, r, cid),
		vMoving: vMoving{
			Delta:    physics.NewVector(0, 0),
			Speed:    physics.NewVector(0, 0),
			Friction: friction,
		},
	}
}

func (iv *InteractVector) Init() event.CID {
	cID := event.NextID(iv)
	iv.CID = cID
	return cID
}

func (iv *InteractVector) ShiftVector(v physics.Vector) {
	iv.Reactive.ShiftPos(v.X(), v.Y())
}
