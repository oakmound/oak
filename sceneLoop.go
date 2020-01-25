package oak

import (
	"image"

	"github.com/oakmound/oak/v2/collision"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/mouse"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
	"github.com/oakmound/oak/v2/timing"
)

var (
	loadingScene = scene.Scene{
		Start: func(prevScene string, data interface{}) {
			dlog.Info("Loading Scene Init")
		},
		Loop: func() bool {
			select {
			case <-startupLoadCh:
				dlog.Info("Load Complete")
				return false
			default:
				return true
			}
		},
		End: func() (string, *scene.Result) {
			return firstScene, nil
		},
	}
)

var firstScene string

func sceneLoop(first string, trackingInputs bool) {
	var prevScene string

	result := new(scene.Result)

	dlog.Info("First Scene Start")

	drawCh <- true
	drawCh <- true

	dlog.Verb("Draw Channel Activated")

	firstScene = first

	SceneMap.CurrentScene = "loading"

	for {
		ViewPos = image.Point{0, 0}
		updateScreen(0, 0)
		useViewBounds = false

		dlog.Info("Scene Start", SceneMap.CurrentScene)
		scen, ok := SceneMap.GetCurrent()
		if !ok {
			dlog.Error("Unknown scene", SceneMap.CurrentScene)
			panic("Unknown scene " + SceneMap.CurrentScene)
		}
		if trackingInputs {
			trackInputChanges()
		}
		go func() {
			dlog.Info("Starting scene in goroutine", SceneMap.CurrentScene)
			scen.Start(prevScene, result.NextSceneInput)
			transitionCh <- true
		}()

		sceneTransition(result)

		// Post transition, begin loading animation
		dlog.Info("Starting load animation")
		drawCh <- true
		dlog.Info("Getting Transition Signal")
		<-transitionCh
		dlog.Info("Resume Drawing")
		// Send a signal to resume (or begin) drawing
		drawCh <- true

		dlog.Info("Looping Scene")
		cont := true

		dlog.ErrorCheck(logicHandler.UpdateLoop(FrameRate, sceneCh))

		for cont {
			select {
			case <-sceneCh:
				cont = scen.Loop()
			case <-skipSceneCh:
				cont = false
			}
		}
		dlog.Info("Scene End", SceneMap.CurrentScene)

		// We don't want enterFrames going off between scenes
		dlog.ErrorCheck(logicHandler.Stop())
		prevScene = SceneMap.CurrentScene

		// Send a signal to stop drawing
		drawCh <- true

		// Reset any ongoing delays
	delayLabel:
		for {
			select {
			case timing.ClearDelayCh <- true:
			default:
				break delayLabel
			}
		}

		dlog.Verb("Resetting Engine")
		// Reset transient portions of the engine
		// We start by clearing the event bus to
		// remove most ongoing code
		logicHandler.Reset()
		// We follow by clearing collision areas
		// because otherwise collision function calls
		// on non-entities (i.e. particles) can still
		// be triggered and attempt to access an entity
		dlog.Verb("Event Bus Reset")
		collision.Clear()
		mouse.Clear()
		event.ResetEntities()
		render.ResetDrawStack()
		render.PreDraw()
		dlog.Verb("Engine Reset")

		// Todo: Add in customizable loading scene between regular scenes,
		// In addition to the existing customizable loading renderable?

		SceneMap.CurrentScene, result = scen.End()
		// For convenience, we allow the user to return nil
		// but it gets translated to an empty result
		if result == nil {
			result = new(scene.Result)
		}

		if !debugResetInProgress {
			debugResetInProgress = true
			go func() {
				debugResetCh <- true
				debugResetInProgress = false
			}()
		}
	}
}
