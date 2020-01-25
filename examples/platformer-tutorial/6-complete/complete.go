package main

import (
	"image/color"
	"math"

	"github.com/oakmound/oak/v2/alg/floatgeom"

	"github.com/oakmound/oak/v2/collision"

	"github.com/oakmound/oak/v2/physics"

	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/key"

	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
)

const (
	// The only collision label we need for this demo is 'ground',
	// indicating something we shouldn't be able to fall or walk through
	Ground collision.Label = 1
)

func main() {
	oak.Add("platformer", func(string, interface{}) {

		char := entities.NewMoving(100, 100, 16, 32,
			render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255}),
			nil, 0, 0)

		render.Draw(char.R)

		char.Speed = physics.NewVector(3, 7)

		fallSpeed := .2

		char.Bind(func(id int, nothing interface{}) int {
			char := event.GetEntity(id).(*entities.Moving)

			// Move left and right with A and D
			if oak.IsDown(key.A) {
				char.Delta.SetX(-char.Speed.X())
			} else if oak.IsDown(key.D) {
				char.Delta.SetX(char.Speed.X())
			} else {
				char.Delta.SetX(0)
			}
			oldX, oldY := char.GetPos()
			char.ShiftPos(char.Delta.X(), char.Delta.Y())

			aboveGround := false

			hit := collision.HitLabel(char.Space, Ground)

			// If we've moved in y value this frame and in the last frame,
			// we were below what we're trying to hit, we are still falling
			if hit != nil && !(oldY != char.Y() && oldY+char.H > hit.Y()) {
				// Correct our y if we started falling into the ground
				char.SetY(hit.Y() - char.H)
				// Stop falling
				char.Delta.SetY(0)
				// Jump with Space when on the ground
				if oak.IsDown(key.Spacebar) {
					char.Delta.ShiftY(-char.Speed.Y())
				}
				aboveGround = true
			} else {
				//Restart when is below ground
				if char.Y() > 500 {
					char.Delta.SetY(0)
					char.SetY(100)
					char.SetX(100)

				}

				// Fall if there's no ground
				char.Delta.ShiftY(fallSpeed)
			}

			if hit != nil {
				// If we walked into a piece of ground, move back
				xover, yover := char.Space.Overlap(hit)
				// We, perhaps unintuitively, need to check the Y overlap, not
				// the x overlap
				// if the y overlap exceeds a superficial value, that suggests
				// we're in a state like
				//
				// G = Ground, C = Character
				//
				// GG C
				// GG C
				//
				// moving to the left
				if math.Abs(yover) > 1 {
					// We add a buffer so this doesn't retrigger immediately
					xbump := 1.0
					if xover > 0 {
						xbump = -1
					}
					char.SetX(oldX + xbump)
					if char.Delta.Y() < 0 {
						char.Delta.SetY(0)
					}
				}

				// If we're below what we hit and we have significant xoverlap, by contrast,
				// then we're about to jump from below into the ground, and we
				// should stop the character.
				if !aboveGround && math.Abs(xover) > 1 {
					// We add a buffer so this doesn't retrigger immediately
					char.SetY(oldY + 1)
					char.Delta.SetY(fallSpeed)
				}

			}

			return 0
		}, event.Enter)

		platforms := []floatgeom.Rect2{
			floatgeom.NewRect2WH(0, 400, 300, 20),
			floatgeom.NewRect2WH(100, 250, 30, 20),
			floatgeom.NewRect2WH(340, 300, 100, 20),
		}

		for _, p := range platforms {
			ground := entities.NewSolid(p.Min.X(), p.Min.Y(), p.W(), p.H(),
				render.NewColorBox(int(p.W()), int(p.H()), color.RGBA{0, 0, 255, 255}),
				nil, 0)
			ground.UpdateLabel(Ground)

			render.Draw(ground.R)
		}

	}, func() bool {
		return true
	}, func() (string, *scene.Result) {
		return "platformer", nil
	})
	oak.Init("platformer")
}
