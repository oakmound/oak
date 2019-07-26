package scene

import (
	"image/color"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
)

func getBasicScene() Scene {
	return Scene{func(string, interface{}) {
		char := entities.NewMoving(100, 100, 16, 32,
			render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255}),
			nil, 0, 0)
		render.Draw(char.R)

	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "basicScene", nil
	},
	}
}  

// Addscene lets a central package manage a set of scenes across subpackages such as in weekly87
// Note the example wont work because there is nothing
func ExampleAddScene() {
	testScene := Scene{}
	AddScene("scene1", getBasicScene())
	AddScene("scene2", getBasicScene())
}
