package main

import (
	"image/color"

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

	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "platformer", nil
	})
	oak.Init("platformer")
}
