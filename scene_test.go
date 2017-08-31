package oak

import (
	"testing"
	"time"

	"github.com/oakmound/oak/scene"
)

func TestSceneTransition(t *testing.T) {
	t.Skip()
	resetOak()
	SetupConfig.Debug = Debug{
		"VERBOSE",
		"",
	}
	AddScene("transition",
		scene.Scene{
			// Initialization function
			Start: func(prevScene string, inData interface{}) {},
			// Loop to continue or stop current scene
			Loop: func() bool { return false },
			// Exit to transition to next scene
			End: func() (nextScene string, result *scene.Result) {
				return "next", &scene.Result{Transition: scene.Fade(.001, 300)}
			},
		})
	Add("next",
		// Initialization function
		func(prevScene string, inData interface{}) {},
		// Loop to continue or stop current scene
		func() bool { return true },
		// Exit to transition to next scene
		func() (nextScene string, result *scene.Result) {
			return "next", nil
		})
	go Init("transition")
	time.Sleep(2 * time.Second)

}
