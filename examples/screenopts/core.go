package main

import (
	"fmt"

	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/key"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
)

func main() {
	oak.Add("demo", func(string, interface{}) {
		txt := render.NewStrText("Press F to toggle fullscreen. Press B to toggle borderless.", 50, 50)
		render.Draw(txt)

		borderless := oak.SetupBorderless
		fullscreen := oak.SetupFullscreen

		event.GlobalBind(func(int, interface{}) int {
			fullscreen = !fullscreen
			err := oak.SetFullScreen(fullscreen)
			if err != nil {
				fullscreen = !fullscreen
				fmt.Println(err)
			}
			return 0
		}, key.Down+key.F)
		event.GlobalBind(func(int, interface{}) int {
			borderless = !borderless
			err := oak.SetBorderless(borderless)
			if err != nil {
				borderless = !borderless
				fmt.Println(err)
			}
			return 0
		}, key.Down+key.B)

	}, func() bool {
		return true
	}, scene.GoTo("demo"))

	// Try uncommenting these
	// Both cannot be true at once!
	// Todo: fix linux bug with client window size not being respected, consuming
	// old border
	// oak.SetupBorderless = true
	// oak.SetupFullscreen = true

	oak.Init("demo")
}
