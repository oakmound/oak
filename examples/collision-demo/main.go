package main

import (
	"image/color"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

const (
	_                   = iota
	RED collision.Label = iota
	GREEN
	BLUE
	TEAL
)

func main() {
	oak.AddScene("demo", scene.Scene{Start: func(ctx *scene.Context) {
		act := &AttachCollisionTest{}
		act.Solid = entities.NewSolid(50, 50, 50, 50, render.NewColorBox(50, 50, color.RGBA{0, 0, 0, 255}), nil, ctx.CallerMap.Register(act))

		collision.Attach(act.Vector, act.Space, nil, 0, 0)

		event.Bind(ctx.EventHandler, event.Enter, act, func(act *AttachCollisionTest, ev event.EnterPayload) event.Response {
			if act.ShouldUpdate {
				act.ShouldUpdate = false
				act.R.Undraw()
				act.R = act.nextR
				render.Draw(act.R, 0, 1)
			}
			if oak.IsDown("A") {
				// We could use attachment here to not have to shift both
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
		})

		render.Draw(act.R, 0, 1)

		collision.PhaseCollision(act.Space, nil)

		event.Bind(ctx.EventHandler, collision.Start, act, func(act *AttachCollisionTest, l collision.Label) event.Response {
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
		})
		event.Bind(ctx.EventHandler, collision.Stop, act, func(act *AttachCollisionTest, l collision.Label) event.Response {
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
		})

		upleft := entities.NewSolid(0, 0, 320, 240, render.NewColorBox(320, 240, color.RGBA{100, 0, 0, 100}), nil, 0)
		upleft.Space.UpdateLabel(RED)
		upleft.R.SetLayer(0)
		render.Draw(upleft.R, 0, 0)

		upright := entities.NewSolid(320, 0, 320, 240, render.NewColorBox(320, 240, color.RGBA{0, 100, 0, 100}), nil, 0)
		upright.Space.UpdateLabel(GREEN)
		upright.R.SetLayer(0)
		render.Draw(upright.R, 0, 0)

		botleft := entities.NewSolid(0, 240, 320, 240, render.NewColorBox(320, 240, color.RGBA{0, 0, 100, 100}), nil, 0)
		botleft.Space.UpdateLabel(BLUE)
		botleft.R.SetLayer(0)
		render.Draw(botleft.R, 0, 0)

		botright := entities.NewSolid(320, 240, 320, 240, render.NewColorBox(320, 240, color.RGBA{0, 100, 100, 100}), nil, 0)
		botright.Space.UpdateLabel(TEAL)
		botright.R.SetLayer(0)
		render.Draw(botright.R, 0, 0)
	}})
	render.SetDrawStack(
		render.NewDynamicHeap(),
	)
	oak.Init("demo")
}

type AttachCollisionTest struct {
	*entities.Solid
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

// CID returns the event.CallerID so that this can be bound to.
func (act *AttachCollisionTest) CID() event.CallerID {
	return act.CallerID
}

// UpdateR with the rgb set on the act.
func (act *AttachCollisionTest) UpdateR() {
	act.nextR = render.NewColorBox(50, 50, color.RGBA{uint8(act.r), uint8(act.g), uint8(act.b), 255})
	act.nextR.SetPos(act.X(), act.Y())
	act.nextR.SetLayer(1)
	act.ShouldUpdate = true
}
