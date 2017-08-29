package oak

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"time"

	"golang.org/x/exp/shiny/screen"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/timing"
)

var (
	sceneMap = map[string]*Scene{
		"loading": {
			false,
			func(prevScene string, data interface{}) {
				dlog.Info("Loading Scene Init")
			},
			func() bool {
				select {
				case <-startupLoadCh:
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

func sceneTransition(result *SceneResult) {
	if result.Transition != nil {
		i := 0
		tx, _ := screenControl.NewTexture(winBuffer.Bounds().Max)
		cont := true
		for cont {
			cont = result.Transition(winBuffer.RGBA(), i)
			drawLoopPublish(tx)
			i++
			time.Sleep(timing.FPSToDuration(FrameRate))
		}
	}
}

type transitionFunction func(*image.RGBA, int) bool

// TransitionFade is a scene transition that fades to black at a given rate for
// a total of frames frames
func TransitionFade(rate float32, frames int) func(*image.RGBA, int) bool {
	rate *= -1
	return func(buf *image.RGBA, frame int) bool {
		if frame > frames {
			return false
		}
		i := float32(frame)
		draw.Draw(buf, buf.Bounds(), render.Brighten(rate*i)(buf), zeroPoint, screen.Src)
		return true
	}
}

// TransitionZoom transitions a scene by zooming in at a relative position of the screen
// at some defined rate. Reasonable values are < .01 for zoomRate.
func TransitionZoom(xPerc, yPerc float64, frames int, zoomRate float64) func(*image.RGBA, int) bool {
	return func(buf *image.RGBA, frame int) bool {
		if frame > frames {
			return false
		}
		z := render.Zoom(xPerc, yPerc, 1+zoomRate*float64(frame))
		draw.Draw(buf, buf.Bounds(), z(buf), zeroPoint, screen.Src)
		return true
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
	NextSceneInput interface{}
	Transition     transitionFunction
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
