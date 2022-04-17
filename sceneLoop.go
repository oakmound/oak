package oak

import (
	"context"

	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/oakerr"
	"github.com/oakmound/oak/v3/scene"
	"github.com/oakmound/oak/v3/timing"
)

// the oak loading scene is a reserved scene
// for preloading assets
const oakLoadingScene = "oak:loading"

func (w *Window) sceneLoop(first string, trackingInputs bool) {
	var prevScene string

	result := new(scene.Result)

	// kick start the draw loop
	w.drawCh <- struct{}{}
	w.drawCh <- struct{}{}

	w.firstScene = first

	w.SceneMap.CurrentScene = oakLoadingScene

	for {
		w.SetViewport(intgeom.Point2{0, 0})
		w.RemoveViewportBounds()

		dlog.Info(dlog.SceneStarting, w.SceneMap.CurrentScene)
		scen, ok := w.SceneMap.GetCurrent()
		if !ok {
			dlog.Error(dlog.UnknownScene, w.SceneMap.CurrentScene)
			if w.ErrorScene != "" {
				w.SceneMap.CurrentScene = w.ErrorScene
				scen, ok = w.SceneMap.GetCurrent()
				if !ok {
					go w.exitWithError(oakerr.NotFound{InputName: "ErrorScene"})
					return
				}
			} else {
				go w.exitWithError(oakerr.NotFound{InputName: "Scene"})
				return
			}
		}
		if trackingInputs {
			w.trackInputChanges()
		}
		gctx, cancel := context.WithCancel(w.ParentContext)
		go func() {
			scen.Start(&scene.Context{
				Context:       gctx,
				PreviousScene: prevScene,
				SceneInput:    result.NextSceneInput,
				DrawStack:     w.DrawStack,
				Handler:       w.eventHandler,
				CallerMap:     w.CallerMap,
				MouseTree:     w.MouseTree,
				CollisionTree: w.CollisionTree,
				Window:        w,
				State:         &w.State,
			})
			w.transitionCh <- struct{}{}
		}()

		w.sceneTransition(result)
		// Post transition, begin loading animation
		w.drawCh <- struct{}{}
		<-w.transitionCh
		// Send a signal to resume (or begin) drawing
		w.drawCh <- struct{}{}

		dlog.Info(dlog.SceneLooping)

		enterCancel := event.EnterLoop(w.eventHandler, timing.FPSToFrameDelay(w.FrameRate))
		nextSceneOverride := ""

		select {
		case <-w.ParentContext.Done():
			w.Quit()
			cancel()
			return
		case <-w.quitCh:
			cancel()
			return
		case nextSceneOverride = <-w.skipSceneCh:
		}
		cancel()
		dlog.Info(dlog.SceneEnding, w.SceneMap.CurrentScene)

		// We don't want enterFrames going off between scenes
		enterCancel()
		prevScene = w.SceneMap.CurrentScene

		// Send a signal to stop drawing
		w.drawCh <- struct{}{}

		// Reset transient portions of the engine
		// We start by clearing the event bus to
		// remove most ongoing code
		w.eventHandler.Reset()
		// We follow by clearing collision areas
		// because otherwise collision function calls
		// on non-entities (i.e. particles) can still
		// be triggered and attempt to access an entity
		w.CollisionTree.Clear()
		w.MouseTree.Clear()
		w.CallerMap.Clear()
		w.eventHandler.SetCallerMap(w.CallerMap)
		w.DrawStack.Clear()
		w.DrawStack.PreDraw()

		// Todo: Add in customizable loading scene between regular scenes,
		// In addition to the existing customizable loading renderable?

		w.SceneMap.CurrentScene, result = scen.End()
		if nextSceneOverride != "" {
			w.SceneMap.CurrentScene = nextSceneOverride
		}
		// For convenience, we allow the user to return nil
		// but it gets translated to an empty result
		if result == nil {
			result = new(scene.Result)
		}
	}
}
