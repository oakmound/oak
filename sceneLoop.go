package oak

import (
	"image"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/timing"
)

func sceneLoop(firstScene string) {
	var prevScene string

	sceneMap[firstScene].active = true
	globalFirstScene = firstScene
	CurrentScene = "loading"

	result := new(SceneResult)

	dlog.Info("First Scene Start")

	drawCh <- true
	drawCh <- true

	dlog.Verb("Draw Channel Activated")

	for {
		ViewPos = image.Point{0, 0}
		updateScreen(0, 0)
		useViewBounds = false

		dlog.Info("Scene Start", CurrentScene)
		go func() {
			dlog.Info("Starting scene in goroutine", CurrentScene)
			sceneMap[CurrentScene].start(prevScene, result.NextSceneInput)
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
		logicTicker := logicLoop()
		for cont {
			select {
			case <-sceneCh:
				cont = sceneMap[CurrentScene].loop()
			case <-skipSceneCh:
				cont = false
			}
		}
		dlog.Info("Scene End", CurrentScene)

		// We don't want enterFrames going off between scenes
		close(logicTicker)
		prevScene = CurrentScene

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
		event.ResetBus()
		// We follow by clearing collision areas
		// because otherwise collision function calls
		// on non-entities (i.e. particles) can still
		// be triggered and attempt to access an entity
		// Todo:
		dlog.Verb("Event Bus Reset")
		collision.Clear()
		mouse.Clear()
		event.ResetEntities()
		render.ResetDrawStack()
		render.PreDraw()
		dlog.Verb("Engine Reset")

		// Todo: Add in customizable loading scene between regular scenes

		CurrentScene, result = sceneMap[CurrentScene].end()
		// For convenience, we allow the user to return nil
		// but it gets translated to an empty result
		if result == nil {
			result = new(SceneResult)
		}

		eb = event.GetBus()
		if !debugResetInProgress {
			debugResetInProgress = true
			go func() {
				debugResetCh <- true
				debugResetInProgress = false
			}()
		}
	}
}
