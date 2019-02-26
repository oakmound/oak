package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/entities/x/move"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/examples/radarDemo/radar"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
)

const (
	_                   = iota
	RED collision.Label = iota
	BLUE
)

func main() {
	oak.Add("demo", func(string, interface{}) {
		act := &AttachCollisionTest{}

		act.Moving = entities.NewMoving(200, 200, 50, 50, render.NewColorBox(50, 50, color.RGBA{0, 0, 0, 255}), nil, act.Init(), 1)
		act.Moving.Speed = physics.NewVector(1, 1)

		collision.Attach(act.Vector, act.Space, 0, 0)

		act.Bind(func(int, interface{}) int {
			move.WASD(act)

			// Normally this should be farmed out on a type such as a character type
			act.ShiftPos(act.Delta.X(), act.Delta.Y())
			act.R.ShiftX(act.Delta.X())
			act.R.ShiftY(act.Delta.Y())
			act.Delta = physics.NewVector(0, 0)
			if act.ShouldUpdate {
				act.ShouldUpdate = false
				act.R.Undraw()
				act.R = act.nextR
				render.Draw(act.R, 0)
			}
			return 0
		}, "EnterFrame")
		render.Draw(act.R, 0)
		act.R.SetLayer(1)

		// Set collision and the functions to update color on collision
		collision.PhaseCollision(act.Space)
		act.Bind(func(id int, label interface{}) int {
			updateOnCollision(act, true, label.(collision.Label))
			return 0
		}, "CollisionStart")
		act.Bind(func(id int, label interface{}) int {
			updateOnCollision(act, false, label.(collision.Label))
			return 0
		}, "CollisionStop")

		// Create two colors to continue collision-demo
		left := entities.NewSolid(0, 0, 320, 480, render.NewColorBox(320, 480, color.RGBA{100, 0, 0, 10}), nil, 0)
		left.Space.UpdateLabel(RED)
		left.R.SetLayer(0)
		//render.Draw(left.R, 0)

		right := entities.NewSolid(320, 0, 320, 480, render.NewColorBox(320, 480, color.RGBA{0, 100, 100, 10}), nil, 0)
		right.Space.UpdateLabel(BLUE)
		right.R.SetLayer(0)
		//render.Draw(right.R, 0)

		// Create the Radar
		center := radar.RadarPoint{act.Xp(), act.Yp()}
		points := make(map[radar.RadarPoint]color.Color)
		r := radar.NewRadar(25, 25, points, center)

		enemy := NewEnemyOnRadar(float64(200))
		enemy.CID.Bind(standardEnemyMove, "EnterFrame")

		r.AddPoint(radar.RadarPoint{enemy.Xp(), enemy.Yp()}, color.RGBA{255, 255, 0, 0})
		render.Draw(enemy.R, 0)
		render.Draw(r, 0)

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

type EnemyOnRadar struct {
	*entities.Moving
}

func (eor *EnemyOnRadar) Init() event.CID {
	return event.NextID(eor)
}
func NewEnemyOnRadar(y float64) *EnemyOnRadar {
	fmt.Println("Sigh")
	eor := new(EnemyOnRadar)
	eor.Moving = entities.NewMoving(50, y, 50, 50, render.NewColorBox(25, 25, color.RGBA{0, 200, 0, 0}), nil, eor.Init(), 0)
	eor.Speed = physics.NewVector(-1*(rand.Float64()*2+1), rand.Float64()*2-1)
	return eor
}

func standardEnemyMove(id int, nothing interface{}) int {

	eor := event.GetEntity(id).(*EnemyOnRadar)
	eor.ShiftX(eor.Speed.X())
	eor.ShiftY(eor.Speed.Y())
	if eor.X() < 0 {
		eor.Speed = physics.NewVector(math.Abs(eor.Speed.X()), (eor.Speed.Y()))
	}
	if eor.X() > 400 {
		eor.Speed = physics.NewVector(-1*math.Abs(eor.Speed.X()), (eor.Speed.Y()))
	}
	if eor.Y() < 0 {
		eor.Speed = physics.NewVector(eor.Speed.X(), math.Abs(eor.Speed.Y()))
	}
	if eor.Y() > 400 {
		eor.Speed = physics.NewVector(eor.Speed.X(), -1*math.Abs(eor.Speed.Y()))
	}
	return 0
}

// updateOnCollision helper function to update color for this small example file
func updateOnCollision(obj *AttachCollisionTest, start bool, label collision.Label) {
	updateValue := 125
	if !start {
		updateValue *= -1
	}
	switch label {
	case RED:
		obj.r += updateValue
	case BLUE:
		obj.b += updateValue

	default:
		return
	}
	obj.UpdateR()
}

type AttachCollisionTest struct {
	*entities.Moving
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
	return event.NextID(act)
}

func (act *AttachCollisionTest) UpdateR() {
	act.nextR = render.NewColorBox(50, 50, color.RGBA{uint8(act.r), uint8(act.g), uint8(act.b), 255})
	act.nextR.SetPos(act.X(), act.Y())
	act.nextR.SetLayer(1)
	act.ShouldUpdate = true
}
