package main

import (
	"image/color"
	"math/rand"

	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

// Rooms exercises shifting the camera in a zelda-esque fashion,
// moving the camera to center on even-sized rooms arranged in a grid
// once the player enters them.

func isOffScreen(ctx *scene.Context, char *entities.Entity) (intgeom.Dir2, bool) {
	x := int(char.X())
	y := int(char.Y())
	if x > ctx.Window.Viewport().X()+ctx.Window.Bounds().X() {
		return intgeom.Right, true
	}
	if y > ctx.Window.Viewport().Y()+ctx.Window.Bounds().Y() {
		return intgeom.Down, true
	}
	if int(char.Right()) < ctx.Window.Viewport().X() {
		return intgeom.Left, true
	}
	if int(char.Bottom()) < ctx.Window.Viewport().Y() {
		return intgeom.Up, true
	}
	return intgeom.Dir2{}, false
}

const (
	transitionFrameCount = 20
)

func main() {

	oak.AddScene("rooms", scene.Scene{Start: func(ctx *scene.Context) {
		char := entities.New(ctx,
			entities.WithRect(floatgeom.NewRect2WH(200, 200, 50, 50)),
			entities.WithColor(color.RGBA{255, 255, 255, 255}),
			entities.WithSpeed(floatgeom.Point2{3, 3}),
			entities.WithDrawLayers([]int{1, 2}),
		)
		var transitioning bool
		var totalTransitionDelta intgeom.Point2
		var transitionDelta intgeom.Point2
		event.Bind(ctx, event.Enter, char, func(c *entities.Entity, ev event.EnterPayload) event.Response {
			dir, ok := isOffScreen(ctx, char)
			if !transitioning && ok {
				transitioning = true
				totalTransitionDelta = ctx.Window.Bounds().Mul(intgeom.Point2{dir.X(), dir.Y()})
				transitionDelta = totalTransitionDelta.DivConst(transitionFrameCount)
			}
			if transitioning {
				// disable movement
				// move camera one size towards the player
				if totalTransitionDelta.X() != 0 || totalTransitionDelta.Y() != 0 {
					oak.ShiftViewport(transitionDelta)
					totalTransitionDelta = totalTransitionDelta.Sub(transitionDelta)
				} else {
					transitioning = false
				}
			} else {
				entities.WASD(char)
			}

			return 0
		})
		for x := 0; x < 2000; x += 12 {
			for y := 0; y < 2000; y += 12 {
				r := uint8(rand.Intn(120))
				g := uint8(rand.Intn(120) + 40)
				cb := render.NewColorBox(12, 12, color.RGBA{r, g, 0, 255})
				cb.SetPos(float64(x), float64(y))
				render.Draw(cb, 0)
			}
		}
	}})

	oak.Init("rooms")
}
