package main

import (
	"image/color"
	"math/rand"

	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/alg/intgeom"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/entities/x/move"
	"github.com/oakmound/oak/v2/physics"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
)

// Rooms exercises shifting the camera in a zelda-esque fashion,
// moving the camera to center on even-sized rooms arranged in a grid
// once the player enters them.

func isOffScreen(char *entities.Moving) (intgeom.Dir2, bool) {
	x := int(char.X())
	y := int(char.Y())
	if x > oak.ViewPos.X+oak.ScreenWidth {
		return intgeom.Right, true
	}
	if y > oak.ViewPos.Y+oak.ScreenHeight {
		return intgeom.Down, true
	}
	if x+int(char.W) < oak.ViewPos.X {
		return intgeom.Left, true
	}
	if y+int(char.H) < oak.ViewPos.Y {
		return intgeom.Up, true
	}
	return intgeom.Dir2{}, false
}

const (
	transitionFrameCount = 20
)

func main() {

	oak.Add("rooms", func(string, interface{}) {
		char := entities.NewMoving(200, 200, 50, 50, render.NewColorBox(50, 50, color.RGBA{125, 125, 0, 255}), nil, 0, 1)
		char.Speed = physics.NewVector(3, 3)

		var transitioning bool
		var totalTransitionDelta intgeom.Point2
		var transitionDelta intgeom.Point2
		char.Bind(func(int, interface{}) int {
			dir, ok := isOffScreen(char)
			if !transitioning && ok {
				transitioning = true
				totalTransitionDelta = intgeom.Point2{oak.ScreenWidth, oak.ScreenHeight}.Mul(intgeom.Point2{dir.X(), dir.Y()})
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
		}, "EnterFrame")
		render.Draw(char.R, 1, 2)

		for x := 0; x < 2000; x += 64 {
			for y := 0; y < 2000; y += 64 {
				r := uint8(rand.Intn(120))
				b := uint8(rand.Intn(120))
				cb := render.NewColorBox(64, 64, color.RGBA{r, 0, b, 255})
				cb.SetPos(float64(x), float64(y))
				render.Draw(cb, 0)
			}
		}

	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "rooms", nil
	})

	oak.Init("rooms")
}
