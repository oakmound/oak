package oak

import (
	"image/color"

	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

// Use oak to display a scene with a single movable character
func Example() {
	AddScene("basicScene", scene.Scene{Start: func(*scene.Context) {
		char := entities.NewMoving(100, 100, 16, 32,
			render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255}),
			nil, 0, 0)
		render.Draw(char.R)
	}})
	Init("basicScene")
}
