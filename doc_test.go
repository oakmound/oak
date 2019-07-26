package oak

import (
	"image/color"

	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
)

// Use oak to display a scene with a single movable character
func Example() { 
	Add("basicScene", func(string, interface{}) {

		char := entities.NewMoving(100, 100, 16, 32,
			render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255}),
			nil, 0, 0)
		render.Draw(char.R)

	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "basicScene", nil
	})
	Init("basicScene")

}

func ExampleAdd() {
	Add("basicScene", func(string, interface{}) { // Whatever you want to do while in the scene
	}, func() bool { // return whether this scene should loop or exit on end
		return true
	}, func() (string, *scene.Result) { // What scene to progress to, make sure its set up!
		return "sceneToBeImplemented", nil
	})
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
 