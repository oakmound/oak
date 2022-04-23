package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"strconv"

	oak "github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

func main() {
	const (
		borderlessAtStart = false
		fullscreenAtStart = false
		topMostAtStart    = false
	)

	oak.AddScene("demo", scene.Scene{Start: func(ctx *scene.Context) {
		txt := render.NewText("Press F to toggle fullscreen. Press B to toggle borderless. Press T to toggle topmost", 50, 50)
		render.Draw(txt)
		line2 := render.NewText("Press Q to change window title. Press C to change the window icon", 50, 70)
		render.Draw(line2)

		borderless := borderlessAtStart
		fullscreen := fullscreenAtStart
		topMost := topMostAtStart

		event.GlobalBind(ctx, key.Down(key.C), func(k key.Event) event.Response {
			colors := []color.RGBA{
				{255, 255, 0, 255},
				{255, 0, 255, 255},
				{0, 255, 255, 255},
				{255, 0, 0, 255},
				{0, 255, 0, 255},
				{0, 0, 255, 255},
			}
			c := colors[rand.Intn(len(colors))]
			rgba := image.NewRGBA(image.Rect(0, 0, 32, 32))
			for x := 0; x < 32; x++ {
				for y := 0; y < 32; y++ {
					rgba.SetRGBA(x, y, c)
				}
			}

			// TODO: this is silly
			err := ctx.Window.(*oak.Window).Window.SetIcon(rgba)
			if err != nil {
				fmt.Println(err)
			}
			return 0
		})

		event.GlobalBind(ctx, key.Down(key.F), func(k key.Event) event.Response {
			fullscreen = !fullscreen
			fmt.Println("Setting fullscreen:", fullscreen)
			err := oak.SetFullScreen(fullscreen)
			if err != nil {
				fullscreen = !fullscreen
				fmt.Println(err)
			}
			return 0
		})
		event.GlobalBind(ctx, key.Down(key.B), func(k key.Event) event.Response {
			borderless = !borderless
			fmt.Println("Setting borderless:", borderless)
			err := oak.SetBorderless(borderless)
			if err != nil {
				borderless = !borderless
				fmt.Println(err)
			}
			return 0
		})
		event.GlobalBind(ctx, key.Down(key.T), func(k key.Event) event.Response {
			topMost = !topMost
			fmt.Println("Setting top most:", topMost)
			err := oak.SetTopMost(topMost)
			if err != nil {
				topMost = !topMost
				fmt.Println(err)
			}
			return 0
		})
		titleCt := 0
		event.GlobalBind(ctx, key.Down(key.Q), func(k key.Event) event.Response {
			titleCt++
			oak.SetTitle("window title " + strconv.Itoa(titleCt))
			return 0
		})

	}})

	oak.Init("demo", func(c oak.Config) (oak.Config, error) {
		c.TopMost = topMostAtStart
		// Both cannot be true at once!
		c.Borderless = borderlessAtStart
		c.Fullscreen = fullscreenAtStart
		return c, nil
	})
}
