package main

import (
	"image/color"

	"github.com/oakmound/oak/v2/physics"

	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/key"

	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
)

func main() {
	oak.Add("platformer", func(string, interface{}) {

		char := entities.NewMoving(100, 100, 16, 32,
			render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255}),
			nil, 0, 0)

		render.Draw(char.R)

		char.Speed = physics.NewVector(3, 3)

		char.Bind(func(id int, nothing interface{}) int {
			char := event.GetEntity(id).(*entities.Moving)
			// Move left and right with A and D
			if oak.IsDown(key.A) {
				char.ShiftX(-char.Speed.X())
			}
			if oak.IsDown(key.D) {
				char.ShiftX(char.Speed.X())
			}
			return 0
		}, event.Enter)

	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "platformer", nil
	})
	oak.Init("platformer")
}
