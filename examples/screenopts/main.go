package main

import (
	"fmt"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

const (
	borderlessAtStart = false
	fullscreenAtStart = false
)

func main() {
	oak.AddScene("demo", scene.Scene{Start: func(*scene.Context) {
		txt := render.NewText("Press F to toggle fullscreen. Press B to toggle borderless.", 50, 50)
		render.Draw(txt)

		borderless := borderlessAtStart
		fullscreen := fullscreenAtStart

		event.GlobalBind(key.Down+key.F, func(event.CID, interface{}) int {
			fullscreen = !fullscreen
			err := oak.SetFullScreen(fullscreen)
			if err != nil {
				fullscreen = !fullscreen
				fmt.Println(err)
			}
			return 0
		})
		event.GlobalBind(key.Down+key.B, func(event.CID, interface{}) int {
			borderless = !borderless
			err := oak.SetBorderless(borderless)
			if err != nil {
				borderless = !borderless
				fmt.Println(err)
			}
			return 0
		})

	}})

	oak.Init("demo", func(c oak.Config) (oak.Config, error) {
		// Both cannot be true at once!
		c.Borderless = borderlessAtStart
		c.Fullscreen = fullscreenAtStart
		return c, nil
	})
}
