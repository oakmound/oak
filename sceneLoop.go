package oak

import (
	"image"

	"bitbucket.org/oakmoundstudio/oak/collision"
	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/mouse"
	"bitbucket.org/oakmoundstudio/oak/render"
	"bitbucket.org/oakmoundstudio/oak/timing"
)

func SceneLoop(firstScene string) {
	var prevScene string

	sceneMap[firstScene].active = true
	globalFirstScene = firstScene
	CurrentScene = "loading"

	result := new(SceneResult)

	dlog.Info("First Scene Start")

	drawChannel <- true
	drawChannel <- true

	for {
		ViewPos = image.Point{0, 0}
		updateScreen(0, 0)
		useViewBounds = false

		dlog.Info("~~~~~~~~~~~Scene Start~~~~~~~~~")
		go func() {
			sceneMap[CurrentScene].start(prevScene, result.NextSceneInput)
			transitionCh <- true
		}()
		sceneTransition(result)
		// Post transition, begin loading animation
		drawChannel <- true
		<-transitionCh
		// Send a signal to resume (or begin) drawing
		drawChannel <- true

		cont := true
		logicTicker := LogicLoop()
		for cont {
			select {
			// The quit channel represents a signal
			// for the engine to stop.
			case <-quitCh:
				return
			case <-sceneCh:
				cont = sceneMap[CurrentScene].loop()
			case <-skipSceneCh:
				cont = false
			}
		}
		dlog.Info("~~~~~~~~Scene End~~~~~~~~~~")

		// We don't want enterFrames going off between scenes
		close(logicTicker)
		prevScene = CurrentScene

		// Send a signal to stop drawing
		drawChannel <- true

		// Reset any ongoing delays
	delayLabel:
		for {
			select {
			case timing.ClearDelayCh <- true:
			default:
				break delayLabel
			}
		}
		// Reset transient portions of the engine
		// We start by clearing the event bus to
		// remove most ongoing code
		event.ResetEventBus()
		// We follow by clearing collision areas
		// because otherwise collision function calls
		// on non-entities (i.e. particles) can still
		// be triggered and attempt to access an entity
		// Todo:
		collision.Clear()
		mouse.Clear()
		event.ResetEntities()
		render.ResetDrawStack()
		render.PreDraw()

		// Todo: Add in customizable loading scene between regular scenes

		CurrentScene, result = sceneMap[CurrentScene].end()
		// For convenience, we allow the user to return nil
		// but it gets translated to an empty result
		if result == nil {
			result = new(SceneResult)
		}

		eb = event.GetEventBus()
		if !debugResetInProgress {
			debugResetInProgress = true
			go func() {
				debugResetCh <- true
				debugResetInProgress = false
			}()
		}
	}
}
