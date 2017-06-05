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
				dlog.Info("Loading Scene Init")
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
// Theses will be refactored to use function pointers sometime
const (
	TRANSITION_NONE = iota
	TRANSITION_FADE
)

func sceneTransition(result *SceneResult) {
	switch result.TransitionType {
	case TRANSITION_FADE:
		tx, err := screenControl.NewTexture(winBuffer.Bounds().Max)
		if err != nil {
			panic(err)
		}
		darkBuffer := winBuffer.RGBA()
		data := result.TransitionPayload.([2]float64)
		rate := float32(data[1]) * -1
		length := float32(data[0])
		for i := float32(0); i < length; i++ {
			draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(),
				render.Brighten(rate*i)(darkBuffer), zeroPoint, screen.Src)
			drawLoopPublish(tx)
		}
	case TRANSITION_NONE:
	default:
	}
}

// A Scene is a set of functions defining what needs to happen when a scene
// starts, loops, and ends.
type Scene struct {
	active bool
	start  SceneStart
	loop   SceneUpdate
	end    SceneEnd
}

// A SceneResult is a set of options for what should be passed into the next
// scene and how the next scene should be transitioned to.
type SceneResult struct {
	NextSceneInput    interface{}
	TransitionType    int
	TransitionPayload interface{}
}

// SceneEnd is a function returning the next scene and a SceneResult
type SceneEnd func() (string, *SceneResult)

// SceneStart is a function taking in a previous scene and a payload
// of data from the previous scene's end
type SceneStart func(prevScene string, data interface{})

// SceneUpdate is a function that returns whether or not the current scene
// should continue to loop.
type SceneUpdate func() bool

// GetScene returns the scene struct with the given name
func GetScene(s string) *Scene {
	return sceneMap[s]
}

// AddScene adds a scene with the given name and functions to the scene map
func AddScene(name string, start SceneStart, loop SceneUpdate, end SceneEnd) error {
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
