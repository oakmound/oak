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

func (c *Controller) sceneLoop(first string, trackingInputs bool) {
	c.SceneMap.AddScene(oakLoadingScene, scene.Scene{
		Start: func(*scene.Context) {
			// TODO: language
			dlog.Info("Loading Scene Init")
		},
		Loop: func() bool {
			return c.startupLoading
		},
		End: func() (string, *scene.Result) {
			dlog.Info("Load Complete")
			return c.firstScene, &scene.Result{
				NextSceneInput: c.FirstSceneInput,
			}
		},
	})

	var prevScene string

	result := new(scene.Result)

	// TODO: language
	dlog.Info("First Scene Start")

	c.drawCh <- struct{}{}
	c.drawCh <- struct{}{}

	// TODO: language
	dlog.Verb("Draw Channel Activated")

	c.firstScene = first

	c.SceneMap.CurrentScene = oakLoadingScene

	for {
		c.setViewport(intgeom.Point2{0, 0})
		c.RemoveViewportBounds()

		dlog.Info("Scene Start: ", c.SceneMap.CurrentScene)
		scen, ok := c.SceneMap.GetCurrent()
		if !ok {
			dlog.Error("Unknown scene: ", c.SceneMap.CurrentScene)
			if c.ErrorScene != "" {
				c.SceneMap.CurrentScene = c.ErrorScene
				scen, ok = c.SceneMap.GetCurrent()
				if !ok {
					go c.exitWithError(oakerr.NotFound{InputName: "ErrorScene"})
					return
				}
			} else {
				go c.exitWithError(oakerr.NotFound{InputName: "Scene"})
				return
			}
		}
		if trackingInputs {
			c.trackInputChanges()
		}
		gctx, cancel := context.WithCancel(context.Background())
		go func() {
			dlog.Info("Starting scene in goroutine", c.SceneMap.CurrentScene)
			scen.Start(&scene.Context{
				Context:       gctx,
				PreviousScene: prevScene,
				SceneInput:    result.NextSceneInput,
				DrawStack:     c.DrawStack,
				EventHandler:  c.logicHandler,
				CallerMap:     c.CallerMap,
				MouseTree:     c.MouseTree,
				CollisionTree: c.CollisionTree,
				Window:        c,
			})
			c.transitionCh <- struct{}{}
		}()

		c.sceneTransition(result)

		// Post transition, begin loading animation
		dlog.Info("Starting load animation")
		c.drawCh <- struct{}{}
		dlog.Info("Getting Transition Signal")
		<-c.transitionCh
		dlog.Info("Resume Drawing")
		// Send a signal to resume (or begin) drawing
		c.drawCh <- struct{}{}

		dlog.Info("Looping Scene")
		cont := true

		dlog.ErrorCheck(c.logicHandler.UpdateLoop(c.FrameRate, c.sceneCh))

		nextSceneOverride := ""

		for cont {
			select {
			case <-c.sceneCh:
				cont = scen.Loop()
			case nextSceneOverride = <-c.skipSceneCh:
				cont = false
			}
		}
		cancel()
		dlog.Info("Scene End", c.SceneMap.CurrentScene)

		// We don't want enterFrames going off between scenes
		dlog.ErrorCheck(c.logicHandler.Stop())
		prevScene = c.SceneMap.CurrentScene

		// Send a signal to stop drawing
		c.drawCh <- struct{}{}

		dlog.Verb("Resetting Engine")
		// Reset transient portions of the engine
		// We start by clearing the event bus to
		// remove most ongoing code
		c.logicHandler.Reset()
		// We follow by clearing collision areas
		// because otherwise collision function calls
		// on non-entities (i.e. particles) can still
		// be triggered and attempt to access an entity
		dlog.Verb("Event Bus Reset")
		c.CollisionTree.Clear()
		c.MouseTree.Clear()
		if c.CallerMap == event.DefaultCallerMap {
			event.ResetCallerMap()
			c.CallerMap = event.DefaultCallerMap
		} else {
			c.CallerMap = event.NewCallerMap()
		}
		c.DrawStack.Clear()
		c.DrawStack.PreDraw()
		dlog.Verb("Engine Reset")

		// Todo: Add in customizable loading scene between regular scenes,
		// In addition to the existing customizable loading renderable?

		c.SceneMap.CurrentScene, result = scen.End()
		if nextSceneOverride != "" {
			c.SceneMap.CurrentScene = nextSceneOverride
		}
		// For convenience, we allow the user to return nil
		// but it gets translated to an empty result
		if result == nil {
			result = new(scene.Result)
		}
	}
}
