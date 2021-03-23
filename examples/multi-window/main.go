package main

import (
	"fmt"
	"image/color"

	"github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
)

func main() {
	oak.SetupConfig.Debug.Level = "VERBOSE"
	// Status 7:00 March 22: creates one window and doesn't error
	c1 := oak.NewController()
	c1.InitialDrawStack = render.NewDrawStack(render.NewStaticHeap())
	// Two windows cannot share the same logic handler
	c1.SetLogicHandler(event.NewBus())
	c1.AddScene("scene1", scene.Scene{
		Start: func(ctx *scene.Context) {
			fmt.Println("Start scene 1")
			cb := render.NewColorBox(50, 50, color.RGBA{255, 0, 0, 255})
			cb.SetPos(50, 50)
			ctx.DrawStack.Draw(cb, 0)
		},
	})
	go c1.Init("scene1")

	c2 := oak.NewController()
	c2.InitialDrawStack = render.NewDrawStack(render.NewStaticHeap())
	c2.SetLogicHandler(event.NewBus())
	c2.AddScene("scene2", scene.Scene{
		Start: func(ctx *scene.Context) {
			fmt.Println("Start scene 2")
			cb := render.NewColorBox(50, 50, color.RGBA{0, 255, 0, 255})
			cb.SetPos(50, 50)
			ctx.DrawStack.Draw(cb, 0)
		},
	})
	c2.Init("scene2")

	//oak.Init() => oak.NewController(render.GlobalDrawStack, dlog.DefaultLogger ...).Init()
}
