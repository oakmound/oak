package oak

import (
	"image/color"

	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
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

// Use AddCommand to grant access to command line commands. Often used to toggle debug modes.
func ExampleAddCommand() {
	debug := true
	AddCommand("SetDebug", func(args []string) {

		if len(args) == 0 {
			debug = !debug
		}
		switch args[0][:1] {
		case "t", "T":
			debug = true
		case "f", "F":
			debug = false
		}

	})
}
