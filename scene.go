package oak

import (
	"errors"
	"fmt"
	"image/draw"

	"golang.org/x/exp/shiny/screen"

	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/render"
)

var (
	sceneMap = map[string]*Scene{
		"loading": {
			false,
			func(prevScene string, data interface{}) {
				return
			},
			func() bool {
				select {
				case <-startupLoadComplete:
					dlog.Info("Load Complete")
					return false
				default:
					return true
				}
			},
			func() (string, *SceneResult) {
				return globalFirstScene, nil
			},
		},
	}
)

// Transition types
const (
	TRANSITION_NONE = iota
	TRANSITION_FADE
)

func sceneTransition(result *SceneResult) {
	switch result.TransitionType {
	case TRANSITION_FADE:
		fmt.Println("Transition Fade starting")
		darkBuffer := winBuffer.RGBA()
		data := result.TransitionPayload.([2]float64)
		rate := float32(data[1]) * -1
		length := float32(data[0])
		for i := float32(0); i < length; i++ {
			draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(),
				render.Brighten(darkBuffer, rate*i), zeroPoint, screen.Src)
			windowControl.Upload(zeroPoint, winBuffer, winBuffer.Bounds())
			windowControl.Publish()
		}
	case TRANSITION_NONE:
	default:
	}
}

type Scene struct {
	active bool
	start  SceneStart
	loop   SceneLoop
	end    SceneEnd
}

type SceneResult struct {
	NextSceneInput    interface{}
	TransitionType    int
	TransitionPayload interface{}
}

type SceneEnd func() (string, *SceneResult)
type SceneStart func(prevScene string, data interface{})
type SceneLoop func() bool

func GetScene(s string) *Scene {
	return sceneMap[s]
}

func AddScene(name string, start SceneStart, loop SceneLoop, end SceneEnd) error {
	fmt.Println("[oak]-------- Adding", name)
	if _, ok := sceneMap[name]; !ok {
		sceneMap[name] = &(Scene{
			false,
			start,
			loop,
			end,
		})
		return nil
	}
	return errors.New("The scene " + name + " is already defined.")
}
