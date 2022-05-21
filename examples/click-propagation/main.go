package main

import (
	"image"
	"image/color"

	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/collision"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/mouse"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

// This example demonstrates the use of the StopPropagation boolean on
// mouse event payloads to prevent mouse interactions from falling
// through to lower UI elements after interacting with a higher layer

func main() {
	oak.AddScene("click-propagation", scene.Scene{
		Start: func(ctx *scene.Context) {
			newHoverButton(ctx, 10, 10, 620, 460, color.RGBA{255, 255, 100, 255}, 1)

			newHoverButton(ctx, 30, 30, 190, 430, color.RGBA{255, 100, 100, 255}, 2)
			newHoverButton(ctx, 240, 30, 370, 430, color.RGBA{255, 100, 255, 255}, 2)

			const gridW = 10
			for x := 50; x < 210-gridW; x += (gridW * 2) {
				for y := 50; y < 450-gridW; y += (gridW * 2) {
					newHoverButton(ctx, float64(x), float64(y), gridW, gridW, color.RGBA{100, 255, 255, 255}, 3)
				}
			}

			newHoverButton(ctx, 260, 50, 100, 390, color.RGBA{100, 100, 255, 255}, 3)
			for y := 70; y < 440-gridW; y += (gridW * 2) {
				newHoverButton(ctx, 270, float64(y), 80, gridW, color.RGBA{255, 255, 255, 255}, 4)
			}
			newHoverButton(ctx, 380, 50, 200, 80, color.RGBA{100, 100, 100, 255}, 3)
		},
	})
	oak.Init("click-propagation")
}

type hoverButton struct {
	id event.CallerID

	mouse.CollisionPhase
	*render.ColorBoxR
}

func (hb *hoverButton) CID() event.CallerID {
	return hb.id
}

func newHoverButton(ctx *scene.Context, x, y, w, h float64, clr color.RGBA, layer int) {
	hb := &hoverButton{}
	hb.id = ctx.Register(hb)
	hb.ColorBoxR = render.NewColorBoxR(int(w), int(h), clr)
	hb.ColorBoxR.SetPos(x, y)

	sp := collision.NewSpace(x, y, w, h, hb.id)
	sp.SetZLayer(float64(layer))

	mouse.Add(sp)
	mouse.PhaseCollision(sp, ctx.Handler)

	render.Draw(hb.ColorBoxR, layer)

	event.Bind(ctx, mouse.ClickOn, hb, func(box *hoverButton, me *mouse.Event) event.Response {
		box.ColorBoxR.Color = image.NewUniform(color.RGBA{128, 128, 128, 128})
		me.StopPropagation = true
		return 0
	})
	event.Bind(ctx, mouse.Start, hb, func(box *hoverButton, me *mouse.Event) event.Response {
		box.ColorBoxR.Color = image.NewUniform(color.RGBA{50, 50, 50, 50})
		me.StopPropagation = true
		return 0
	})
	event.Bind(ctx, mouse.Stop, hb, func(box *hoverButton, me *mouse.Event) event.Response {
		box.ColorBoxR.Color = image.NewUniform(clr)
		me.StopPropagation = true
		return 0
	})
}
