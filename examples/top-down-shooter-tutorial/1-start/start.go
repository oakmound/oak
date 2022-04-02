package main

import (
	"image/color"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
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

		event.Bind(ctx, event.Enter, char, func(c *entities.Moving, ev event.EnterPayload) event.Response {

			char.Delta.Zero()
			if oak.IsDown(key.WStr) {
				char.Delta.ShiftY(-char.Speed.Y())
			}
			if oak.IsDown(key.AStr) {
				char.Delta.ShiftX(-char.Speed.X())
			}
			if oak.IsDown(key.SStr) {
				char.Delta.ShiftY(char.Speed.Y())
			}
			if oak.IsDown(key.DStr) {
				char.Delta.ShiftX(char.Speed.X())
			}
			char.ShiftPos(char.Delta.X(), char.Delta.Y())
			hit := char.HitLabel(Enemy)
			if hit != nil {
				playerAlive = false
			}

			return 0
		})

	}})
	oak.Init("tds")
}
