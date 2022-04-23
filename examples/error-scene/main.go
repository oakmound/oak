package main

import (
	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

func main() {
	controller := oak.NewWindow()
	// If ErrorScene is set, the scene handler will
	// fall back to this error scene if it is told to
	// go to an unknown scene
	controller.ErrorScene = "error"
	controller.AddScene("typo", scene.Scene{Start: func(ctx *scene.Context) {
		ctx.DrawStack.Draw(render.NewText("Real scene", 100, 100))
	}})
	controller.AddScene("error", scene.Scene{Start: func(ctx *scene.Context) {
		ctx.DrawStack.Draw(render.NewText("Error scene", 100, 100))
	}})

	controller.Init("typpo")
}
