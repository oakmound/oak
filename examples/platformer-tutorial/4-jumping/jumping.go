package main

import (
	"image/color"

	"github.com/oakmound/oak/v3/collision"

	"github.com/oakmound/oak/v3/physics"

	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

// Collision labels
const (
	// The only collision label we need for this demo is 'ground',
	// indicating something we shouldn't be able to fall or walk through
	Ground collision.Label = 1
)

func main() {
	oak.AddScene("platformer", scene.Scene{Start: func(ctx *scene.Context) {

		char := entities.NewMoving(100, 100, 16, 32,
			render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255}),
			nil, 0, 0)

		render.Draw(char.R)

		char.Speed = physics.NewVector(3, 3)

		fallSpeed := .1

		event.Bind(ctx, event.Enter, char, func(c *entities.Moving, ev event.EnterPayload) event.Response {
			// Move left and right with A and D
			if oak.IsDown(key.A) {
				c.ShiftX(-c.Speed.X())
			}
			if oak.IsDown(key.D) {
				c.ShiftX(c.Speed.X())
			}
			hit := collision.HitLabel(c.Space, Ground)
			if hit == nil {
				// Fall if there's no ground
				c.Delta.ShiftY(fallSpeed)
			} else {
				c.Delta.SetY(0)
				// Jump with Space
				if oak.IsDown(key.Spacebar) {
					c.Delta.ShiftY(-c.Speed.Y())
				}
			}
			c.ShiftY(c.Delta.Y())
			return 0
		})

		ground := entities.NewSolid(0, 400, 500, 20,
			render.NewColorBox(500, 20, color.RGBA{0, 0, 255, 255}),
			nil, 0)
		ground.UpdateLabel(Ground)

		render.Draw(ground.R)

	}})
	oak.Init("platformer")
}
