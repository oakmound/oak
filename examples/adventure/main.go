package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/alg/intgeom"
	"github.com/oakmound/oak/v4/collision"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
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
	transitionFrameCount = 25
)

const (
	LabelPlayer collision.Label = iota
	LabelWall
	LabelBox
	LabelLock
	LabelKey
	LabelPot
	LabelTreasure
	LabelEnemy
)

func main() {

	oak.AddScene("rooms", scene.Scene{Start: func(ctx *scene.Context) {
		char := entities.New(ctx,
			entities.WithRect(floatgeom.NewRect2WH(200, 200, 50, 50)),
			entities.WithColor(color.RGBA{255, 255, 255, 255}),
			entities.WithSpeed(floatgeom.Point2{5, 5}),
			entities.WithDrawLayers([]int{1, 2}),
		)
		var transitioning bool
		var totalTransitionDelta intgeom.Point2
		var transitionDelta intgeom.Point2
		event.Bind(ctx, event.Enter, char, func(c *entities.Entity, ev event.EnterPayload) event.Response {
			if !transitioning {
				dir, ok := isOffScreen(ctx, char)
				if ok {
					transitioning = true
					totalTransitionDelta = ctx.Window.Bounds().Mul(intgeom.Point2{dir.X(), dir.Y()})
					transitionDelta = totalTransitionDelta.DivConst(transitionFrameCount)
				}
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
				char.Delta = floatgeom.Point2{}
				if ctx.IsDown(key.W) {
					char.Delta[1] -= char.Speed[1]
				}
				if ctx.IsDown(key.S) {
					char.Delta[1] += char.Speed[1]
				}
				if ctx.IsDown(key.A) {
					char.Delta[0] -= char.Speed[0]
				}
				if ctx.IsDown(key.D) {
					char.Delta[0] += char.Speed[0]
				}
				char.Rect = char.Rect.Shift(char.Delta)

				char.Tree.UpdateSpace(
					char.X(), char.Y(), char.W(), char.H(), char.Space,
				)
				hitWall := false
				hits := char.Tree.Hits(char.Space)
				for _, h := range hits {
					fmt.Println("hit label", h.Label)
					switch h.Label {
					case LabelWall:
						if hitWall {
							continue
						}
						hitWall = true
						char.Delta = char.Delta.MulConst(-1)
						char.Rect = char.Rect.Shift(char.Delta)
						char.Tree.UpdateSpace(
							char.X(), char.Y(), char.W(), char.H(), char.Space,
						)
					}
				}
				char.Renderable.SetPos(char.X(), char.Y())
			}

			return 0
		})
		const tileWidth = 50
		const tileHeight = 50
		tileDims := entities.WithDimensions(floatgeom.Point2{tileWidth, tileHeight})
		x := 0.0
		y := float64(-tileHeight) // to accommodate for initial newline
		for _, rn := range board {
			if rn == '\n' {
				x = 0
				y += tileHeight
				continue
			}
			commonOpts := entities.And(
				tileDims,
				entities.WithPosition(floatgeom.Point2{x, y}),
				entities.WithDrawLayers([]int{1}),
			)
			switch Tile(rn) {
			case Wall:
				entities.New(ctx, commonOpts,
					entities.WithColor(color.RGBA{50, 50, 50, 255}),
					entities.WithLabel(LabelWall),
				)
				x += tileWidth
				continue
			case Box:
				entities.New(ctx, commonOpts,
					tileDims,
					entities.WithColor(color.RGBA{150, 150, 20, 255}),
					entities.WithLabel(LabelBox),
				)
				// TODO
			case Player:
				fmt.Println("placing character at", x, y)
				char.SetPos(floatgeom.Point2{x, y})
			case Pot:
				entities.New(ctx, commonOpts,
					tileDims,
					entities.WithColor(color.RGBA{100, 100, 100, 255}),
					entities.WithLabel(LabelTreasure),
				)
				// TODO
			case Enemy:
				// TODO
			case Treasure:
				entities.New(ctx, commonOpts,
					tileDims,
					entities.WithColor(color.RGBA{255, 255, 0, 255}),
					entities.WithLabel(LabelTreasure),
				)
			case Lock:
				entities.New(ctx, commonOpts,
					tileDims,
					entities.WithColor(color.RGBA{200, 150, 0, 255}),
					entities.WithLabel(LabelLock),
				)
			case Key:
				entities.New(ctx, commonOpts,
					tileDims,
					entities.WithColor(color.RGBA{50, 50, 200, 255}),
					entities.WithLabel(LabelKey),
				)
			case Empty:
			}

			// ground beneath moving objects
			r := uint8(rand.Intn(20) + 40)
			g := uint8(rand.Intn(80) + 60)
			b := uint8(rand.Intn(10) + 30)
			cb := render.NewColorBoxR(tileWidth, tileHeight, color.RGBA{r, g, b, 255})
			cb.SetPos(float64(x), float64(y))
			render.Draw(cb, 0)

			x += tileWidth
		}
	}})

	oak.Init("rooms", func(c oak.Config) (oak.Config, error) {
		c.Screen.Width = 650
		c.Screen.Height = 500
		return c, nil
	})
}

type Tile rune

const (
	Wall     Tile = 'W'
	Empty    Tile = ' '
	Box      Tile = 'B'
	Enemy    Tile = 'E'
	Player   Tile = 'C'
	Pot      Tile = 'P'
	Lock     Tile = 'L'
	Key      Tile = 'K'
	Treasure Tile = 'T'
)

const board = `
WWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWW
W           WW                        WW           W
W     C     WW                        WW           W
W           WW                        WWW          W
W                                     LLL        T W
W                                     WWW          W
W           WW                        WW           W
W           WW                        WW           W
W           WW                     P  WWWWWWWWWWWWWW
WWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWW    WW           W
WWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWW    WW       B   W
W                        WW           WW      BKB  W
W                        WW           WW  E    B   W
W       WWWWWWWWWWW      WW    E      WW           W
W           WW           WW           WW      E    W
W           WW                        WW           W
W           WW                        WW           W
W           WW                                     W
W           WW                                     W
WWWWWW      WWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWW
WWWWWW      WWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWW
W                         B  B     P  WW           W
W  P                      BB          WW    K      W
W                         B  B    K   WW      E    W
W                         BB       P  WW           W
W           WWWWWWWWWWWWWWWWWWWWWWWWWWWW           W
W                                     B        E   W
W      E                              P    E       W
W                                     B            W
WWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWW
`
