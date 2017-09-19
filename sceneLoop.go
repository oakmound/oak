package oak

import (
	"image"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/oak/timing"
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
			return "", nil
		},
	}
)

func sceneLoop(firstScene string) {
	var prevScene string

	result := new(scene.Result)

	dlog.Info("First Scene Start")

	drawCh <- true
	drawCh <- true

	dlog.Verb("Draw Channel Activated")

	loadingScene.End = func() (string, *scene.Result) {
		return firstScene, nil
	}

	// Todo: consider changing the name of this scene to avoid collisions
	err := SceneMap.AddScene("loading", loadingScene)
	if err != nil {
		dlog.Error("Loading scene unable to be added", err)
		panic("Loading scene unable to be added")
	}
	SceneMap.CurrentScene = "loading"

	for {
		ViewPos = image.Point{0, 0}
		updateScreen(0, 0)
		useViewBounds = false

		dlog.Info("Scene Start", SceneMap.CurrentScene)
		scen, ok := SceneMap.GetCurrent()
		if !ok {
			dlog.Error("Unknown scene", SceneMap.CurrentScene)
			panic("Unknown scene")
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

		err = logicHandler.UpdateLoop(FrameRate, sceneCh)
		if err != nil {
			dlog.Error(err)
		}

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
		err = logicHandler.Stop()
		if err != nil {
			dlog.Error(err)
		}
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
