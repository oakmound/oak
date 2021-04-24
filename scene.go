package oak

import (
	"time"

	"github.com/oakmound/oak/v2/scene"
	"github.com/oakmound/oak/v2/timing"
)

// AddScene is shorthand for oak.SceneMap.AddScene
func (c *Controller) AddScene(name string, s scene.Scene) error {
	return c.SceneMap.AddScene(name, s)
}

func (c *Controller) sceneTransition(result *scene.Result) {
	if result.Transition != nil {
		i := 0
		tx, _ := c.screenControl.NewTexture(c.winBuffer.Bounds().Max)
		cont := true
		for cont {
			cont = result.Transition(c.winBuffer.RGBA(), i)
			c.drawLoopPublish(c, tx)
			i++
			time.Sleep(timing.FPSToFrameDelay(c.FrameRate))
		}
	}
}
