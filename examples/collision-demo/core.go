package main

import (
	"image/color"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
)

const (
	_                   = iota
	RED collision.Label = iota
	GREEN
	BLUE
	TEAL
)

func main() {
	oak.Add("demo", func(string, interface{}) {
		act := &AttachCollisionTest{}
		act.Solid = entities.NewSolid(50, 50, 50, 50, render.NewColorBox(50, 50, color.RGBA{0, 0, 0, 255}), nil, act.Init())

		collision.Attach(act.Vector, act.Space, 0, 0)

		act.Bind(func(int, interface{}) int {
			if act.ShouldUpdate {
				act.ShouldUpdate = false
				act.R.Undraw()
				act.R = act.nextR
				render.Draw(act.R, 0)
			}
			if oak.IsDown("A") {
				// We could use attachement here to not have to shift both
				// R and act but that is made more difficult by constantly
				// changing the act's R
				act.ShiftX(-3)
				act.R.ShiftX(-3)
			} else if oak.IsDown("D") {
				act.ShiftX(3)
				act.R.ShiftX(3)
			}
			if oak.IsDown("W") {
				act.ShiftY(-3)
				act.R.ShiftY(-3)
			} else if oak.IsDown("S") {
				act.ShiftY(3)
				act.R.ShiftY(3)
			}
			return 0
		}, "EnterFrame")

		render.Draw(act.R, 0)
		act.R.SetLayer(1)

		collision.PhaseCollision(act.Space)
		act.Bind(func(id int, label interface{}) int {
			l := label.(collision.Label)
			switch l {
			case RED:
				act.r += 125
				act.UpdateR()
			case GREEN:
				act.g += 125
				act.UpdateR()
			case BLUE:
				act.b += 125
				act.UpdateR()
			case TEAL:
				act.b += 125
				act.g += 125
				act.UpdateR()
			}
			return 0
		}, "CollisionStart")
		act.Bind(func(id int, label interface{}) int {
			l := label.(collision.Label)
			switch l {
			case RED:
				act.r -= 125
				act.UpdateR()
			case GREEN:
				act.g -= 125
				act.UpdateR()
			case BLUE:
				act.b -= 125
				act.UpdateR()
			case TEAL:
				act.b -= 125
				act.g -= 125
				act.UpdateR()
			}
			return 0
		}, "CollisionStop")

		upleft := entities.NewSolid(0, 0, 320, 240, render.NewColorBox(320, 240, color.RGBA{100, 0, 0, 100}), nil, 0)
		upleft.Space.UpdateLabel(RED)
		upleft.R.SetLayer(0)
		render.Draw(upleft.R, 0)

		upright := entities.NewSolid(320, 0, 320, 240, render.NewColorBox(320, 240, color.RGBA{0, 100, 0, 100}), nil, 0)
		upright.Space.UpdateLabel(GREEN)
		upright.R.SetLayer(0)
		render.Draw(upright.R, 0)

		botleft := entities.NewSolid(0, 240, 320, 240, render.NewColorBox(320, 240, color.RGBA{0, 0, 100, 100}), nil, 0)
		botleft.Space.UpdateLabel(BLUE)
		botleft.R.SetLayer(0)
		render.Draw(botleft.R, 0)

		botright := entities.NewSolid(320, 240, 320, 240, render.NewColorBox(320, 240, color.RGBA{0, 100, 100, 100}), nil, 0)
		botright.Space.UpdateLabel(TEAL)
		botright.R.SetLayer(0)
		render.Draw(botright.R, 0)

	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "demo", nil
	})
	render.SetDrawStack(
		render.NewHeap(false),
		render.NewDrawFPS(),
	)
	oak.Init("demo")
}

type AttachCollisionTest struct {
	entities.Solid
	// AttachSpace is a composable struct that allows
	// spaces to be attached to vectors
	collision.AttachSpace
	// Phase is a composable struct that enables the call
	// collision.CollisionPhase on this struct's space,
	// which will start sending signals when that space
	// starts and stops touching given labels
	collision.Phase
	r, g, b      int
	ShouldUpdate bool
	nextR        render.Renderable
}

func (act *AttachCollisionTest) Init() event.CID {
	cid := event.NextID(act)
	act.CID = cid
	return cid
}

func (act *AttachCollisionTest) UpdateR() {
	act.nextR = render.NewColorBox(50, 50, color.RGBA{uint8(act.r), uint8(act.g), uint8(act.b), 255})
	act.nextR.SetPos(act.X(), act.Y())
	act.nextR.SetLayer(1)
	act.ShouldUpdate = true
}
