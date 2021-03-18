package oak

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v2/scene"
)

func TestSceneTransition(t *testing.T) {
	t.Skip()
	resetOak()
	SetupConfig.Debug = Debug{
		"VERBOSE",
		"",
	}
	err := Add("next",
		// Initialization function
		func(*scene.Context) {},
		// Loop to continue or stop current scene
		func() bool { return true },
		// Exit to transition to next scene
		func() (nextScene string, result *scene.Result) {
			return "next", nil
		})
	if err != nil {
		t.Fatalf("add scene failed: %v", err)
	}
	go Init("transition")
	time.Sleep(2 * time.Second)

}
