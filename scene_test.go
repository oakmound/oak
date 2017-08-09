package oak

import (
	"testing"
	"time"
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
	AddScene("transition",
		// Initialization function
		func(prevScene string, inData interface{}) {},
		// Loop to continue or stop current scene
		func() bool { return false },
		// Exit to transition to next scene
		func() (nextScene string, result *SceneResult) {
			return "next", &SceneResult{Transition: TransitionFade(.001, 300)}
		})
	AddScene("next",
		// Initialization function
		func(prevScene string, inData interface{}) {},
		// Loop to continue or stop current scene
		func() bool { return true },
		// Exit to transition to next scene
		func() (nextScene string, result *SceneResult) {
			return "next", nil
		})
	go Init("transition")
	time.Sleep(2 * time.Second)

}
