package oak

import (
	"testing"
	"time"

	"github.com/oakmound/oak/scene"
)

func TestBadScene(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	sceneLoop("badscene")
}

func TestSceneTransition(t *testing.T) {
	t.Skip()
	resetOak()
	SetupConfig.Debug = Debug{
		"VERBOSE",
		"",
	}
	SceneMap.Add("transition",
		// Initialization function
		func(prevScene string, inData interface{}) {},
		// Loop to continue or stop current scene
		func() bool { return false },
		// Exit to transition to next scene
		func() (nextScene string, result *scene.Result) {
			return "next", &scene.Result{Transition: scene.Fade(.001, 300)}
		})
	SceneMap.Add("next",
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
