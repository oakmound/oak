package main

import (
	"image/color"

	"github.com/oakmound/oak/v3/physics"

	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

func main() {
	oak.AddScene("platformer", scene.Scene{Start: func(ctx *scene.Context) {

		char := entities.NewMoving(100, 100, 16, 32,
			render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255}),
			nil, 0, 0)

		render.Draw(char.R)

		char.Speed = physics.NewVector(3, 3)
		event.Bind(ctx, event.Enter, char, func(c *entities.Moving, ev event.EnterPayload) event.Response {
			// Move left and right with A and D
			if oak.IsDown(key.A) {
				c.ShiftX(-c.Speed.X())
			}
			if oak.IsDown(key.D) {
				c.ShiftX(c.Speed.X())
			}
			return 0
		})
	}})
	oak.Init("platformer")
}
