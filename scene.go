package oak

import (
	"time"

	"github.com/oakmound/oak/v3/scene"
	"github.com/oakmound/oak/v3/timing"
)

// AddScene is shorthand for c.SceneMap.AddScene
func (w *Window) AddScene(name string, s scene.Scene) error {
	return w.SceneMap.AddScene(name, s)
}

func (w *Window) sceneTransition(result *scene.Result) {
	if result.Transition != nil {
		i := 0
		cont := true
		frameDelay := timing.FPSToFrameDelay(w.FrameRate)
		for cont {
			cont = result.Transition(w.winBuffer.RGBA(), i)
			w.publish()
			i++
			time.Sleep(frameDelay)
		}
	}
}
