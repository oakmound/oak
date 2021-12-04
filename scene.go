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
		frameDelay := timing.FPSToFrameDelay(w.DrawFrameRate)
		for cont {
			// TODO: Transition should take in the amount of time passed, not number of frames,
			// to account for however long the transition itself takes.
			cont = result.Transition(w.winBuffers[w.bufferIdx].RGBA(), i)
			w.publish()
			i++
			time.Sleep(frameDelay)
		}
	}
}
