package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/oakmound/oak/v3/render/mod"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/collision/ray"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/physics"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

// Collision labels
const (
	Enemy  collision.Label = 1
	Player collision.Label = 2
)

var (
	// Vectors are backed by pointers,
	// so despite this not being a pointer,
	// this does update according to the player's
	// position so long as we don't reset
	// the player's position vector
	playerPos physics.Vector

	destroy = event.RegisterEvent[event.NoPayload]()

	sheet [][]*render.Sprite
)

func main() {
	oak.AddScene("tds", scene.Scene{Start: func(ctx *scene.Context) {
		// Initialization

		sprites, err := render.GetSheet("sheet.png")
		dlog.ErrorCheck(err)
		sheet = sprites.ToSprites()

		// Player setup
		eggplant, err := render.GetSprite("eggplant-fish.png")
		playerR := render.NewSwitch("left", map[string]render.Modifiable{
			"left": eggplant,
			// We must copy the sprite before we modify it, or "left"
			// will also be flipped.
			"right": eggplant.Copy().Modify(mod.FlipX),
		})
		if err != nil {
			fmt.Println(err)
		}
		char := entities.NewMoving(100, 100, 32, 32,
			playerR,
			nil, 0, 0)

		char.Speed = physics.NewVector(5, 5)
		playerPos = char.Point.Vector
		render.Draw(char.R, 2)

		event.Bind(ctx, event.Enter, char, func(char *entities.Moving, ev event.EnterPayload) event.Response {
			char.Delta.Zero()
			if oak.IsDown(key.WStr) {
				char.Delta.ShiftY(-char.Speed.Y())
			}
			if oak.IsDown(key.AStr) {
				char.Delta.ShiftX(-char.Speed.X())
			}
			if oak.IsDown(key.SStr) {
				char.Delta.ShiftY(char.Speed.Y())
			}
			if oak.IsDown(key.DStr) {
				char.Delta.ShiftX(char.Speed.X())
			}
			char.ShiftPos(char.Delta.X(), char.Delta.Y())
			hit := char.HitLabel(Enemy)
			if hit != nil {
				ctx.Window.NextScene()
			}

			// update animation
			swtch := char.R.(*render.Switch)
			if char.Delta.X() > 0 {
				if swtch.Get() == "left" {
					swtch.Set("right")
				}
			} else if char.Delta.X() < 0 {
				if swtch.Get() == "right" {
					swtch.Set("left")
				}
			}

			return 0
		})

		event.Bind(ctx, mouse.Press, char, func(char *entities.Moving, mevent *mouse.Event) event.Response {
			x := char.X() + char.W/2
			y := char.Y() + char.H/2
			ray.DefaultCaster.CastDistance = floatgeom.Point2{x, y}.Sub(floatgeom.Point2{mevent.X(), mevent.Y()}).Magnitude()
			hits := ray.CastTo(floatgeom.Point2{x, y}, floatgeom.Point2{mevent.X(), mevent.Y()})
			for _, hit := range hits {
				event.TriggerForCallerOn(ctx, hit.Zone.CID, destroy, event.NoPayload{})
			}
			ctx.DrawForTime(
				render.NewLine(x, y, mevent.X(), mevent.Y(), color.RGBA{0, 128, 0, 128}),
				time.Millisecond*50,
				2)
			return 0
		})

		// Create enemies periodically
		event.GlobalBind(ctx, event.Enter, func(enterPayload event.EnterPayload) event.Response {
			if enterPayload.FramesElapsed%EnemyRefresh == 0 {
				go NewEnemy(ctx)
			}
			return 0
		})

		// Draw the background
		for x := 0; x < ctx.Window.Width(); x += 16 {
			for y := 0; y < ctx.Window.Height(); y += 16 {
				i := rand.Intn(3) + 1
				// Get a random tile to draw in this position
				sp := sheet[i/2][i%2].Copy()
				sp.SetPos(float64(x), float64(y))
				render.Draw(sp, 1)
			}
		}

	}})

	oak.Init("tds", func(c oak.Config) (oak.Config, error) {
		// This indicates to oak to automatically open and load image and audio
		// files local to the project before starting any scene.
		c.BatchLoad = true
		c.Debug.Level = "Verbose"
		c.Assets.ImagePath = "assets/images"

		return c, nil
	})
}

// Top down shooter constsv
const (
	EnemyRefresh = 30
	EnemySpeed   = 2
)

// NewEnemy creates an enemy for a top down shooter
func NewEnemy(ctx *scene.Context) {
	x, y := enemyPos(ctx)

	enemyFrame := sheet[0][0].Copy()
	enemyR := render.NewSwitch("left", map[string]render.Modifiable{
		"left":  enemyFrame,
		"right": enemyFrame.Copy().Modify(mod.FlipX),
	})
	enemy := entities.NewSolid(x, y, 16, 16,
		enemyR,
		nil, 0)

	render.Draw(enemy.R, 2)

	enemy.UpdateLabel(Enemy)

	event.Bind(ctx, event.Enter, enemy, func(e *entities.Solid, ev event.EnterPayload) event.Response {
		// move towards the player
		x, y := enemy.GetPos()
		pt := floatgeom.Point2{x, y}
		pt2 := floatgeom.Point2{playerPos.X(), playerPos.Y()}
		delta := pt2.Sub(pt).Normalize().MulConst(EnemySpeed)
		enemy.ShiftPos(delta.X(), delta.Y())

		// update animation
		swtch := enemy.R.(*render.Switch)
		if delta.X() > 0 {
			if swtch.Get() == "left" {
				swtch.Set("right")
			}
		} else if delta.X() < 0 {
			if swtch.Get() == "right" {
				swtch.Set("left")
			}
		}
		return 0
	})

	event.Bind(ctx, destroy, enemy, func(e *entities.Solid, nothing event.NoPayload) event.Response {
		e.Destroy()
		return 0
	})
}

func enemyPos(ctx *scene.Context) (float64, float64) {
	w := ctx.Window.Width()
	h := ctx.Window.Height()
	// Spawn on the edge of the screen
	perimeter := w*2 + h*2
	pos := int(rand.Float64() * float64(perimeter))
	// Top
	if pos < w {
		return float64(pos), 0
	}
	pos -= w
	// Right
	if pos < h {
		return float64(w), float64(pos)
	}
	// Bottom
	pos -= h
	if pos < w {
		return float64(pos), float64(h)
	}
	pos -= w
	// Left
	return 0, float64(pos)
}
