package oak

import (
	"time"

	"github.com/oakmound/oak/v3/scene"
	"github.com/oakmound/oak/v3/timing"
)

// AddScene is shorthand for c.SceneMap.AddScene
func (c *Controller) AddScene(name string, s scene.Scene) error {
	return c.SceneMap.AddScene(name, s)
}

func (c *Controller) sceneTransition(result *scene.Result) {
	if result.Transition != nil {
		i := 0
		cont := true
		for cont {
			cont = result.Transition(c.winBuffer.RGBA(), i)
			c.publish()
			i++
			time.Sleep(timing.FPSToFrameDelay(c.FrameRate))
		}
	}
}
