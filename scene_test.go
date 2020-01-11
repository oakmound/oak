package oak

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v2/scene"
	"github.com/stretchr/testify/assert"
)

func TestSceneTransition(t *testing.T) {
	t.Skip()
	resetOak()
	SetupConfig.Debug = Debug{
		"VERBOSE",
		"",
	}
	err := AddScene("transition",
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
	assert.Nil(t, err)
	err = Add("next",
		// Initialization function
		func(prevScene string, inData interface{}) {},
		// Loop to continue or stop current scene
		func() bool { return true },
		// Exit to transition to next scene
		func() (nextScene string, result *scene.Result) {
			return "next", nil
		})
	assert.Nil(t, err)
	go Init("transition")
	time.Sleep(2 * time.Second)

}
