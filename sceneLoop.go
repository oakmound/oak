package oak

import (
	"context"
	"fmt"

	"github.com/oakmound/oak/v2/alg/intgeom"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
	"github.com/oakmound/oak/v2/timing"
)

func (c *Controller) sceneLoop(first string, trackingInputs bool, debugConsoleDisabled bool) {
	err := c.SceneMap.AddScene("loading", scene.Scene{
		Start: func(*scene.Context) {
			// TODO: language
			dlog.Info("Loading Scene Init")
		},
		Loop: func() bool {
			select {
			case <-c.startupLoadCh:
				// TODO: language
				dlog.Info("Load Complete")
				return false
			default:
				fmt.Println("loading still")
				return true
			}
		},
		End: func() (string, *scene.Result) {
			return c.firstScene, &scene.Result{
				NextSceneInput: c.FirstSceneInput,
			}
		},
	})
	if err != nil {
		// ???
	}

	var prevScene string

	result := new(scene.Result)

	// TODO: language
	dlog.Info("First Scene Start")

	c.drawCh <- struct{}{}
	c.drawCh <- struct{}{}

	// TODO: language
	dlog.Verb("Draw Channel Activated")

	c.firstScene = first

	c.SceneMap.CurrentScene = "loading"

	for {
		c.ViewPos = intgeom.Point2{0, 0}
		c.updateScreen(c.ViewPos)
		c.useViewBounds = false

		dlog.Info("Scene Start", c.SceneMap.CurrentScene)
		scen, ok := c.SceneMap.GetCurrent()
		if !ok {
			dlog.Error("Unknown scene", c.SceneMap.CurrentScene)
			if c.ErrorScene != "" {
				c.SceneMap.CurrentScene = c.ErrorScene
				scen, ok = c.SceneMap.GetCurrent()
				if !ok {
					panic("error scene not defined in scene map")
				}
			} else {
				panic("Unknown scene " + c.SceneMap.CurrentScene)
			}
		}
		if trackingInputs {
			trackInputChanges()
		}
		gctx, cancel := context.WithCancel(context.Background())
		defer cancel()
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

		for cont {
			select {
			case <-c.sceneCh:
				cont = scen.Loop()
			case <-c.skipSceneCh:
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

		// Reset any ongoing delays
	delayLabel:
		for {
			select {
			case timing.ClearDelayCh <- struct{}{}:
			default:
				break delayLabel
			}
		}

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
		if c.DrawStack == render.GlobalDrawStack {
			render.ResetDrawStack()
			c.DrawStack = render.GlobalDrawStack
		} else {
			c.DrawStack = c.InitialDrawStack.Copy()
		}
		c.DrawStack.PreDraw()
		dlog.Verb("Engine Reset")

		// Todo: Add in customizable loading scene between regular scenes,
		// In addition to the existing customizable loading renderable?

		c.SceneMap.CurrentScene, result = scen.End()
		// For convenience, we allow the user to return nil
		// but it gets translated to an empty result
		if result == nil {
			result = new(scene.Result)
		}

		if !debugConsoleDisabled && !c.debugResetInProgress {
			c.debugResetInProgress = true
			go func() {
				c.debugResetCh <- struct{}{}
				c.debugResetInProgress = false
			}()
		}
	}
}
