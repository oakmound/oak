package main

import (
	"image/color"
	"math/rand"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/entities/x/move"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/physics"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

// Rooms exercises shifting the camera in a zelda-esque fashion,
// moving the camera to center on even-sized rooms arranged in a grid
// once the player enters them.

func isOffScreen(ctx *scene.Context, char *entities.Moving) (intgeom.Dir2, bool) {
	x := int(char.X())
	y := int(char.Y())
	if x > ctx.Window.Viewport().X()+ctx.Window.Width() {
		return intgeom.Right, true
	}
	if y > ctx.Window.Viewport().Y()+ctx.Window.Height() {
		return intgeom.Down, true
	}
	if x+int(char.W) < ctx.Window.Viewport().X() {
		return intgeom.Left, true
	}
	if y+int(char.H) < ctx.Window.Viewport().Y() {
		return intgeom.Up, true
	}
	return intgeom.Dir2{}, false
}

const (
	transitionFrameCount = 20
)

func main() {

	oak.AddScene("rooms", scene.Scene{Start: func(ctx *scene.Context) {
		char := entities.NewMoving(200, 200, 50, 50, render.NewColorBox(50, 50, color.RGBA{255, 255, 255, 255}), nil, 0, 1)
		char.Speed = physics.NewVector(3, 3)

		var transitioning bool
		var totalTransitionDelta intgeom.Point2
		var transitionDelta intgeom.Point2
		char.Bind(event.Enter, func(event.CallerID, interface{}) int {
			dir, ok := isOffScreen(ctx, char)
			if !transitioning && ok {
				transitioning = true
				totalTransitionDelta = intgeom.Point2{ctx.Window.Width(), ctx.Window.Height()}.Mul(intgeom.Point2{dir.X(), dir.Y()})
				transitionDelta = totalTransitionDelta.DivConst(transitionFrameCount)
			}
			if transitioning {
				// disable movement
				// move camera one size towards the player
				if totalTransitionDelta.X() != 0 || totalTransitionDelta.Y() != 0 {
					oak.ShiftScreen(transitionDelta.X(), transitionDelta.Y())
					totalTransitionDelta = totalTransitionDelta.Sub(transitionDelta)
				} else {
					transitioning = false
				}
			} else {
				move.WASD(char)
			}

			return 0
		})
		render.Draw(char.R, 1, 2)

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
