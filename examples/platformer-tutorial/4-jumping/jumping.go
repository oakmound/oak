package main

import (
	"image/color"

	"github.com/oakmound/oak/v2/collision"

	"github.com/oakmound/oak/v2/physics"

	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/key"

	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
)

// Collision labels
const (
	// The only collision label we need for this demo is 'ground',
	// indicating something we shouldn't be able to fall or walk through
	Ground collision.Label = 1
)

func main() {
	oak.Add("platformer", func(string, interface{}) {

		char := entities.NewMoving(100, 100, 16, 32,
			render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255}),
			nil, 0, 0)

		render.Draw(char.R)

		char.Speed = physics.NewVector(3, 3)

		fallSpeed := .1

		char.Bind(func(id int, nothing interface{}) int {
			char := event.GetEntity(id).(*entities.Moving)
			// Move left and right with A and D
			if oak.IsDown(key.A) {
				char.ShiftX(-char.Speed.X())
			}
			if oak.IsDown(key.D) {
				char.ShiftX(char.Speed.X())
			}
			hit := collision.HitLabel(char.Space, Ground)
			if hit == nil {
				// Fall if there's no ground
				char.Delta.ShiftY(fallSpeed)
			} else {
				char.Delta.SetY(0)
				// Jump with Space
				if oak.IsDown(key.Spacebar) {
					char.Delta.ShiftY(-char.Speed.Y())
				}
			}
			char.ShiftY(char.Delta.Y())
			return 0
		}, event.Enter)

		ground := entities.NewSolid(0, 400, 500, 20,
			render.NewColorBox(500, 20, color.RGBA{0, 0, 255, 255}),
			nil, 0)
		ground.UpdateLabel(Ground)

		render.Draw(ground.R)

	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "platformer", nil
	})
	oak.Init("platformer")
}
