package main

import (
	"fmt"
	"image/color"

	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

func main() {
	c1 := oak.NewController()
	c1.DrawStack = render.NewDrawStack(render.NewDynamicHeap())

	// Two windows cannot share the same logic handler
	c1.SetLogicHandler(event.NewBus(nil))
	c1.FirstSceneInput = color.RGBA{255, 0, 0, 255}
	c1.AddScene("scene1", scene.Scene{
		Start: func(ctx *scene.Context) {
			fmt.Println("Start scene 1")
			cb := render.NewColorBox(50, 50, ctx.SceneInput.(color.RGBA))
			cb.SetPos(50, 50)
			ctx.DrawStack.Draw(cb, 0)
			dFPS := render.NewDrawFPS(0.1, nil, 600, 10)
			ctx.DrawStack.Draw(dFPS, 1)
			ctx.EventHandler.GlobalBind(mouse.Press, mouse.Binding(func(_ event.CID, me mouse.Event) int {
				cb.SetPos(me.X(), me.Y())
				return 0
			}))
		},
	})
	go func() {
		c1.Init("scene1", func(c oak.Config) (oak.Config, error) {
			c.Debug.Level = "VERBOSE"
			c.DrawFrameRate = 1200
			c.FrameRate = 60
			return c, nil
		})
		fmt.Println("scene 1 exited")
	}()

	c2 := oak.NewController()
	c2.DrawStack = render.NewDrawStack(render.NewDynamicHeap())
	c2.SetLogicHandler(event.NewBus(nil))
	c2.FirstSceneInput = color.RGBA{0, 255, 0, 255}
	c2.AddScene("scene2", scene.Scene{
		Start: func(ctx *scene.Context) {
			fmt.Println("Start scene 2")
			cb := render.NewColorBox(50, 50, ctx.SceneInput.(color.RGBA))
			cb.SetPos(50, 50)
			ctx.DrawStack.Draw(cb, 0)
			dFPS := render.NewDrawFPS(0.1, nil, 600, 10)
			ctx.DrawStack.Draw(dFPS, 1)
			ctx.EventHandler.GlobalBind(mouse.Press, mouse.Binding(func(_ event.CID, me mouse.Event) int {
				cb.SetPos(me.X(), me.Y())
				return 0
			}))
		},
	})
	c2.Init("scene2", func(c oak.Config) (oak.Config, error) {
		c.Debug.Level = "VERBOSE"
		c.DrawFrameRate = 1200
		c.FrameRate = 60
		return c, nil
	})
	fmt.Println("scene 2 exited")

	//oak.Init() => oak.NewController(render.GlobalDrawStack, dlog.DefaultLogger ...).Init()
}
