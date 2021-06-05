package main

import (
	"image/color"
	"time"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/collision"
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
	playerAlive = true
)

func main() {
	oak.AddScene("tds", scene.Scene{Start: func(ctx *scene.Context) {
		playerAlive = true
		char := entities.NewMoving(100, 100, 32, 32,
			render.NewColorBox(32, 32, color.RGBA{0, 255, 0, 255}),
			nil, 0, 0)

		char.Speed = physics.NewVector(5, 5)
		render.Draw(char.R)

		char.Bind(event.Enter, func(id event.CID, _ interface{}) int {
			char := event.GetEntity(id).(*entities.Moving)
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
				playerAlive = false
			}

			return 0
		})

		char.Bind(mouse.Press, func(id event.CID, me interface{}) int {
			char := event.GetEntity(id).(*entities.Moving)
			mevent := me.(*mouse.Event)
			ctx.DrawForTime(
				render.NewLine(char.X()+char.W/2, char.Y()+char.H/2, mevent.X(), mevent.Y(), color.RGBA{0, 128, 0, 128}),
				time.Millisecond*50,
				1)
			return 0
		})

	}, Loop: func() bool {
		return playerAlive
	}})
	oak.Init("tds")
}
