package main

import (
	"image/color"
	"math/rand"
	"time"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/collision/ray"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/physics"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

// Collision labels
const (
	Enemy  collision.Label = 1
	Player collision.Label = 2
)

var (
	// Vectors are backed by pointers,
	// so despite this not being a pointer,
	// this does update according to the player's
	// position so long as we don't reset
	// the player's position vector
	playerPos physics.Vector

	destroy = event.RegisterEvent[event.NoPayload]()
)

func main() {
	oak.AddScene("tds", scene.Scene{Start: func(ctx *scene.Context) {
		char := entities.NewMoving(100, 100, 32, 32,
			render.NewColorBox(32, 32, color.RGBA{0, 255, 0, 255}),
			nil, 0, 0)

		char.Speed = physics.NewVector(5, 5)
		playerPos = char.Point.Vector
		render.Draw(char.R)

		event.Bind(ctx, event.Enter, char, func(char *entities.Moving, ev event.EnterPayload) event.Response {

			char.Delta.Zero()
			if oak.IsDown(key.W) {
				char.Delta.ShiftY(-char.Speed.Y())
			}
			if oak.IsDown(key.A) {
				char.Delta.ShiftX(-char.Speed.X())
			}
			if oak.IsDown(key.S) {
				char.Delta.ShiftY(char.Speed.Y())
			}
			if oak.IsDown(key.D) {
				char.Delta.ShiftX(char.Speed.X())
			}
			char.ShiftPos(char.Delta.X(), char.Delta.Y())
			hit := char.HitLabel(Enemy)
			if hit != nil {
				ctx.Window.NextScene()
			}
			return 0
		})

		event.Bind(ctx, mouse.Press, char, func(char *entities.Moving, mevent *mouse.Event) event.Response {
			x := char.X() + char.W/2
			y := char.Y() + char.H/2
			ray.DefaultCaster.CastDistance = floatgeom.Point2{x, y}.Sub(floatgeom.Point2{mevent.X(), mevent.Y()}).Magnitude()
			hits := ray.CastTo(floatgeom.Point2{x, y}, floatgeom.Point2{mevent.X(), mevent.Y()})
			for _, hit := range hits {
				event.TriggerForCallerOn(ctx, hit.Zone.CID, destroy, event.NoPayload{})
			}
			ctx.DrawForTime(
				render.NewLine(x, y, mevent.X(), mevent.Y(), color.RGBA{0, 128, 0, 128}),
				time.Millisecond*50,
				1)
			return 0
		})

		event.GlobalBind(ctx, event.Enter, func(enterPayload event.EnterPayload) event.Response {
			if enterPayload.FramesElapsed%EnemyRefresh == 0 {
				go NewEnemy(ctx)
			}
			return 0
		})

	}})
	oak.Init("tds")
}

// Top down shooter consts
const (
	EnemyRefresh = 30
	EnemySpeed   = 2
)

// NewEnemy creates an enemy for a top down shooter
func NewEnemy(ctx *scene.Context) {
	x, y := enemyPos(ctx)

	enemy := entities.NewSolid(x, y, 16, 16,
		render.NewColorBox(16, 16, color.RGBA{200, 0, 0, 200}),
		nil, 0)

	render.Draw(enemy.R)

	enemy.UpdateLabel(Enemy)
	event.Bind(ctx, event.Enter, enemy, func(e *entities.Solid, ev event.EnterPayload) event.Response {
		// move towards the player
		x, y := e.GetPos()
		pt := floatgeom.Point2{x, y}
		pt2 := floatgeom.Point2{playerPos.X(), playerPos.Y()}
		delta := pt2.Sub(pt).Normalize().MulConst(EnemySpeed)
		e.ShiftPos(delta.X(), delta.Y())
		return 0
	})
	event.Bind(ctx, destroy, enemy, func(e *entities.Solid, nothing event.NoPayload) event.Response {
		e.Destroy()
		return 0
	})
}

func enemyPos(ctx *scene.Context) (float64, float64) {
	w := ctx.Window.Width()
	h := ctx.Window.Height()
	// Spawn on the edge of the screen
	perimeter := w*2 + h*2
	pos := int(rand.Float64() * float64(perimeter))
	// Top
	if pos < w {
		return float64(pos), 0
	}
	pos -= w
	// Right
	if pos < h {
		return float64(w), float64(pos)
	}
	// Bottom
	pos -= h
	if pos < w {
		return float64(pos), float64(h)
	}
	pos -= w
	// Left
	return 0, float64(pos)
}
