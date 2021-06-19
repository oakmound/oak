package oak

import (
	"testing"

	"github.com/oakmound/oak/v3/scene"
)

func TestSceneTransition(t *testing.T) {
	c1 := NewWindow()
	c1.AddScene("1", scene.Scene{
		Start: func(context *scene.Context) {
			go context.Window.NextScene()
		},
		End: func() (nextScene string, result *scene.Result) {
			return "2", &scene.Result{
				Transition: scene.Fade(1, 10),
			}
		},
	})
	c1.AddScene("2", scene.Scene{
		Start: func(context *scene.Context) {
			context.Window.Quit()
		},
	})
	c1.Init("1")
}
