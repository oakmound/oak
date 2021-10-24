package oak

import (
	"context"

	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/oakerr"
	"github.com/oakmound/oak/v3/scene"
)

// the oak loading scene is a reserved scene
// for preloading assets
const oakLoadingScene = "oak:loading"

func (w *Window) sceneLoop(first string, trackingInputs bool) {
	w.SceneMap.AddScene(oakLoadingScene, scene.Scene{
		Loop: func() bool {
			return w.startupLoading
		},
		End: func() (string, *scene.Result) {
			return w.firstScene, &scene.Result{
				NextSceneInput: w.FirstSceneInput,
			}
		},
	})

	var prevScene string

	result := new(scene.Result)

	w.drawCh <- struct{}{}
	w.drawCh <- struct{}{}

	w.firstScene = first

	w.SceneMap.CurrentScene = oakLoadingScene

	for {
		w.setViewport(intgeom.Point2{0, 0})
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
				EventHandler:  w.eventHandler,
				CallerMap:     w.CallerMap,
				MouseTree:     w.MouseTree,
				CollisionTree: w.CollisionTree,
				Window:        w,
				KeyState:      &w.State,
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
		cont := true

		dlog.ErrorCheck(w.eventHandler.UpdateLoop(w.FrameRate, w.sceneCh))

		nextSceneOverride := ""

		for cont {
			select {
			case <-w.ParentContext.Done():
			case <-w.quitCh:
				cancel()
				return
			case <-w.sceneCh:
				cont = scen.Loop()
			case nextSceneOverride = <-w.skipSceneCh:
				cont = false
			}
		}
		cancel()
		dlog.Info(dlog.SceneEnding, w.SceneMap.CurrentScene)

		// We don't want enterFrames going off between scenes
		dlog.ErrorCheck(w.eventHandler.Stop())
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
		if w.CallerMap == event.DefaultCallerMap {
			event.ResetCallerMap()
			w.CallerMap = event.DefaultCallerMap
		} else {
			w.CallerMap = event.NewCallerMap()
		}
		if cmr, ok := w.eventHandler.(event.CallerMapper); ok {
			cmr.SetCallerMap(w.CallerMap)
		}
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
