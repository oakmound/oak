package main

import (
	"embed"
	"image/color"
	"math/rand"
	"time"

	"github.com/oakmound/oak/v3/render/mod"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/collision/ray"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

const (
	Enemy collision.Label = 1
)

var (
	playerX *float64
	playerY *float64

	destroy = event.RegisterEvent[struct{}]()

	sheet [][]*render.Sprite
)

const (
	fieldWidth  = 1000
	fieldHeight = 1000
)

func main() {

	oak.AddScene("tds", scene.Scene{Start: func(ctx *scene.Context) {
		render.Draw(render.NewDrawFPS(0, nil, 10, 10), 2, 0)
		render.Draw(render.NewLogicFPS(0, nil, 10, 20), 2, 0)
		// render.Draw(debugtools.NewThickRTree(ctx, collision.DefaultTree, 5), 2, 3)

		// Initialization
		sprites, err := render.GetSheet("sheet.png")
		dlog.ErrorCheck(err)
		sheet = sprites.ToSprites()

		oak.SetViewportBounds(intgeom.NewRect2(0, 0, fieldWidth, fieldHeight))

		// Player setup
		eggplant, err := render.GetSprite("eggplant-fish.png")
		dlog.ErrorCheck(err)
		playerR := render.NewSwitch("left", map[string]render.Modifiable{
			"left": eggplant,
			// We must copy the sprite before we modify it, or "left"
			// will also be flipped.
			"right": eggplant.Copy().Modify(mod.FlipX),
		})
		char := entities.New(ctx,
			entities.WithRect(floatgeom.NewRect2WH(100, 100, 32, 32)),
			entities.WithRenderable(playerR),
			entities.WithSpeed(floatgeom.Point2{3, 3}),
			entities.WithDrawLayers([]int{1, 2}),
		)

		playerX = &char.Rect.Min[0]
		playerY = &char.Rect.Min[1]

		screenCenter := ctx.Window.Bounds().DivConst(2)

		event.Bind(ctx, event.Enter, char, func(char *entities.Entity, ev event.EnterPayload) event.Response {
			if oak.IsDown(key.W) {
				char.Delta[1] += (-char.Speed.Y() * ev.TickPercent)
			}
			if oak.IsDown(key.A) {
				char.Delta[0] += (-char.Speed.X() * ev.TickPercent)
			}
			if oak.IsDown(key.S) {
				char.Delta[1] += (char.Speed.Y() * ev.TickPercent)
			}
			if oak.IsDown(key.D) {
				char.Delta[0] += (char.Speed.X() * ev.TickPercent)
			}
			ctx.Window.(*oak.Window).DoBetweenDraws(func() {
				char.ShiftDelta()
				oak.SetViewport(
					screenCenter.Sub(intgeom.Point2{
						int(char.X()), int(char.Y()),
					}),
				)
				char.Delta = floatgeom.Point2{}
			})
			// Don't go out of bounds
			if char.X() < 0 {
				char.SetX(0)
			} else if char.X() > fieldWidth-char.W() {
				char.SetX(fieldWidth - char.W())
			}
			if char.Y() < 0 {
				char.SetY(0)
			} else if char.Y() > fieldHeight-char.H() {
				char.SetY(fieldHeight - char.H())
			}

			hit := char.HitLabel(Enemy)
			if hit != nil {
				ctx.Window.NextScene()
			}

			// update animation
			swtch := char.Renderable.(*render.Switch)
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

		event.Bind(ctx, mouse.Press, char, func(char *entities.Entity, mevent *mouse.Event) event.Response {
			x := char.X() + char.W()/2
			y := char.Y() + char.H()/2
			vp := ctx.Window.Viewport()
			mx := mevent.X() + float64(vp.X())
			my := mevent.Y() + float64(vp.Y())
			ray.DefaultCaster.CastDistance = floatgeom.Point2{x, y}.Sub(floatgeom.Point2{mx, my}).Magnitude()
			hits := ray.CastTo(floatgeom.Point2{x, y}, floatgeom.Point2{mx, my})
			for _, hit := range hits {
				event.TriggerForCallerOn(ctx, hit.Zone.CID, destroy, struct{}{})
			}
			ctx.DrawForTime(
				render.NewLine(x, y, mx, my, color.RGBA{0, 128, 0, 128}),
				time.Millisecond*50,
				1, 2)
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
		for x := 0; x < fieldWidth; x += 16 {
			for y := 0; y < fieldHeight; y += 16 {
				i := rand.Intn(3) + 1
				// Get a random tile to draw in this position
				sp := sheet[i/2][i%2].Copy()
				sp.SetPos(float64(x), float64(y))
				render.Draw(sp, 0, 1)
			}
		}

	}})

	render.SetDrawStack(
		render.NewCompositeR(),
		render.NewDynamicHeap(),
		render.NewStaticHeap(),
	)

	oak.SetFS(assets)
	oak.Init("tds", func(c oak.Config) (oak.Config, error) {
		c.BatchLoad = true
		c.Assets.ImagePath = "assets/images"
		//c.FrameRate = 30
		return c, nil
	})
}

//go:embed assets
var assets embed.FS

// Top down shooter consts
const (
	EnemyRefresh = 25
	EnemySpeed   = 2
)

// NewEnemy creates an enemy for a top down shooter
func NewEnemy(ctx *scene.Context) {
	x, y := enemyPos()

	enemyFrame := sheet[0][0].Copy()
	enemyR := render.NewSwitch("left", map[string]render.Modifiable{
		"left":  enemyFrame,
		"right": enemyFrame.Copy().Modify(mod.FlipX),
	})
	enemy := entities.New(ctx,
		entities.WithRect(floatgeom.NewRect2WH(x, y, 16, 16)),
		entities.WithRenderable(enemyR),
		entities.WithDrawLayers([]int{1, 2}),
		entities.WithLabel(Enemy),
	)

	event.Bind(ctx, event.Enter, enemy, func(e *entities.Entity, ev event.EnterPayload) event.Response {
		// move towards the player
		x, y := enemy.X(), enemy.Y()
		pt := floatgeom.Point2{x, y}
		pt2 := floatgeom.Point2{*playerX, *playerY}
		delta := pt2.Sub(pt).Normalize().MulConst(EnemySpeed * ev.TickPercent)
		enemy.Shift(delta)

		// update animation
		swtch := enemy.Renderable.(*render.Switch)
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

	event.Bind(ctx, destroy, enemy, func(e *entities.Entity, nothing struct{}) event.Response {
		e.Destroy()
		return 0
	})
}

func enemyPos() (float64, float64) {
	// Spawn on the edge of the screen
	perimeter := fieldWidth*2 + fieldHeight*2
	pos := int(rand.Float64() * float64(perimeter))
	// Top
	if pos < fieldWidth {
		return float64(pos), 0
	}
	pos -= fieldWidth
	// Right
	if pos < fieldHeight {
		return float64(fieldWidth), float64(pos)
	}
	// Bottom
	pos -= fieldHeight
	if pos < fieldWidth {
		return float64(pos), float64(fieldHeight)
	}
	pos -= fieldWidth
	// Left
	return 0, float64(pos)
}
