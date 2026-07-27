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
	"github.com/oakmound/oak/v4/mouse"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
	"github.com/oakmound/oak/v4/window"
)

func main() {
	const (
		borderlessAtStart = false
		fullscreenAtStart = false
		topMostAtStart    = false
	)

	oak.AddScene("demo", scene.Scene{Start: func(ctx *scene.Context) {
		txt := render.NewText("Press F to toggle fullscreen. Press B to toggle borderless. Press T to toggle topmost / floating.", 50, 50)
		render.Draw(txt)
		line2 := render.NewText("Press Q to change window title. Press C to change the window icon. Press H to replace the cursor.", 50, 70)
		render.Draw(line2)
		line3 := render.NewText("Press M to maximize the window. Press Z to minimize the window. Press N to normalize the window.", 50, 90)
		render.Draw(line3)
		line4 := render.NewText("Press P to print the position of the window.", 50, 110)
		render.Draw(line4)

		borderless := borderlessAtStart
		fullscreen := fullscreenAtStart
		topMost := topMostAtStart

		event.GlobalBind(ctx, key.Down(key.C), func(k key.Event) event.Response {
			fmt.Println("Setting icon")
			colors := []color.NRGBA{
				{255, 255, 0, 255},
				{255, 0, 255, 255},
				{0, 255, 255, 255},
				{255, 0, 0, 255},
				{0, 255, 0, 255},
				{0, 0, 255, 255},
			}
			c := colors[rand.Intn(len(colors))]
			rgba := image.NewNRGBA(image.Rect(0, 0, 48, 48))
			for x := range 48 {
				for y := range 48 {
					rgba.SetNRGBA(x, y, c)
				}
			}

			err := ctx.Window.SetIcon(rgba)
			if err != nil {
				fmt.Println(err)
			}
			return 0
		})

		event.GlobalBind(ctx, key.Down(key.F), func(k key.Event) event.Response {
			fullscreen = !fullscreen
			fmt.Println("Setting fullscreen:", fullscreen)
			err := ctx.Window.SetFullScreen(fullscreen)
			if err != nil {
				fullscreen = !fullscreen
				fmt.Println(err)
			}
			return 0
		})
		event.GlobalBind(ctx, key.Down(key.B), func(k key.Event) event.Response {
			borderless = !borderless
			fmt.Println("Setting borderless:", borderless)
			err := ctx.Window.SetBorderless(borderless)
			if err != nil {
				borderless = !borderless
				fmt.Println(err)
			}
			return 0
		})
		event.GlobalBind(ctx, key.Down(key.T), func(k key.Event) event.Response {
			topMost = !topMost
			fmt.Println("Setting top most:", topMost)
			err := ctx.Window.SetTopMost(topMost)
			if err != nil {
				topMost = !topMost
				fmt.Println(err)
			}
			return 0
		})
		event.GlobalBind(ctx, key.Down(key.M), func(k key.Event) event.Response {
			fmt.Println("Maximizing")
			err := ctx.Window.(window.ExtendedWindow).Maximize()
			if err != nil {
				fmt.Println(err)
			}
			return 0
		})
		event.GlobalBind(ctx, key.Down(key.Z), func(k key.Event) event.Response {
			fmt.Println("Minimizing")
			err := ctx.Window.(window.ExtendedWindow).Minimize()
			if err != nil {
				fmt.Println(err)
			}
			return 0
		})
		event.GlobalBind(ctx, key.Down(key.N), func(k key.Event) event.Response {
			fmt.Println("Normalizing")
			err := ctx.Window.(window.ExtendedWindow).Normalize()
			if err != nil {
				fmt.Println(err)
			}
			return 0
		})
		event.GlobalBind(ctx, key.Down(key.P), func(k key.Event) event.Response {
			x, y := ctx.Window.(window.ExtendedWindow).GetDesktopPosition()
			fmt.Println("Position:", x, y)
			return 0
		})
		titleCt := 0
		event.GlobalBind(ctx, key.Down(key.Q), func(k key.Event) event.Response {
			titleCt++
			ctx.Window.SetTitle("window title " + strconv.Itoa(titleCt))
			return 0
		})
		event.GlobalBind(ctx, key.Down(key.H), func(k key.Event) event.Response {
			ctx.Window.HideCursor()
			box := render.NewSequence(15,
				render.NewColorBox(2, 2, color.RGBA{255, 255, 0, 255}),
				render.NewColorBox(3, 3, color.RGBA{255, 235, 0, 255}),
				render.NewColorBox(4, 4, color.RGBA{255, 215, 0, 255}),
				render.NewColorBox(5, 5, color.RGBA{255, 195, 0, 255}),
				render.NewColorBox(6, 6, color.RGBA{255, 175, 0, 255}),
				render.NewColorBox(5, 5, color.RGBA{255, 155, 0, 255}),
				render.NewColorBox(4, 4, color.RGBA{255, 135, 0, 255}),
				render.NewColorBox(3, 3, color.RGBA{255, 115, 0, 255}),
				render.NewColorBox(2, 2, color.RGBA{255, 95, 0, 255}),
				render.NewColorBox(1, 1, color.RGBA{255, 75, 0, 255}),
				render.EmptyRenderable(),
				render.EmptyRenderable(),
				render.EmptyRenderable(),
				render.EmptyRenderable(),
			)
			ctx.DrawStack.Draw(box)

			event.GlobalBind(ctx,
				mouse.Drag, func(mouseEvent *mouse.Event) event.Response {
					box.SetPos(mouseEvent.X(), mouseEvent.Y())
					return 0
				})
			return event.ResponseUnbindThisBinding
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
