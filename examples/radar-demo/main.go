package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/alg/intgeom"

	oak "github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/examples/radar-demo/radar"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

const (
	xLimit = 1000
	yLimit = 1000
)

// This example demonstrates making a basic radar or other custom renderable
// type. The radar here acts as a UI element, staying on screen, and follows
// around a player character.
//TODO: Remove and or link to grove radar as it is cleaner
// https://github.com/oakmound/grove/tree/master/components/radar

func main() {
	oak.AddScene("demo", scene.Scene{Start: func(ctx *scene.Context) {
		render.Draw(render.NewDrawFPS(0.03, nil, 10, 10))

		char := entities.New(ctx,
			entities.WithRect(floatgeom.NewRect2WH(200, 200, 50, 50)),
			entities.WithColor(color.RGBA{125, 125, 0, 255}),
			entities.WithSpeed(floatgeom.Point2{3, 3}),
			entities.WithDrawLayers([]int{1, 2}),
		)

		oak.SetViewportBounds(intgeom.NewRect2(0, 0, xLimit, yLimit))
		moveRect := floatgeom.NewRect2(0, 0, xLimit, yLimit)
		event.Bind(ctx, event.Enter, char, func(char *entities.Entity, ev event.EnterPayload) event.Response {
			entities.WASD(char)
			entities.Limit(char, moveRect)
			entities.CenterScreenOn(char)
			return 0
		})

		// Create the Radar
		xp := &char.Rect.Min[0]
		yp := &char.Rect.Min[1]
		center := radar.Point{X: xp, Y: yp}
		points := make(map[radar.Point]color.Color)
		w := 100
		h := 100
		r := radar.NewRadar(w, h, points, center, 10)
		r.SetPos(float64(ctx.Window.Bounds().X()-w), 0)

		for i := 0; i < 5; i++ {
			x, y := rand.Float64()*400, rand.Float64()*400
			enemy := newEnemyOnRadar(ctx, x, y)
			event.Bind(ctx, event.Enter, enemy, standardEnemyMove)
			xp := &enemy.Rect.Min[0]
			yp := &enemy.Rect.Min[1]
			r.AddPoint(radar.Point{X: xp, Y: yp}, color.RGBA{255, 255, 0, 255})
		}

		render.Draw(r, 2)

		for x := 0; x < xLimit; x += 64 {
			for y := 0; y < yLimit; y += 64 {
				r := uint8(rand.Intn(120))
				b := uint8(rand.Intn(120))
				cb := render.NewColorBox(64, 64, color.RGBA{r, 0, b, 255})
				cb.SetPos(float64(x), float64(y))
				render.Draw(cb, 0)
			}
		}

	}})

	render.SetDrawStack(
		render.NewCompositeR(),
		render.NewDynamicHeap(),
		render.NewStaticHeap(),
	)
	oak.Init("demo")
}

func newEnemyOnRadar(ctx *scene.Context, x, y float64) *entities.Entity {
	eor := entities.New(ctx,
		entities.WithRect(floatgeom.NewRect2WH(50, y, 50, 50)),
		entities.WithColor(color.RGBA{0, 200, 0, 200}),
		entities.WithSpeed(floatgeom.Point2{-1 * (rand.Float64()*2 + 1), rand.Float64()*2 - 1}),
		entities.WithDrawLayers([]int{1, 1}),
	)
	eor.Delta = eor.Speed
	return eor
}

func standardEnemyMove(eor *entities.Entity, ev event.EnterPayload) event.Response {
	if eor.X() < 0 {
		eor.Delta = floatgeom.Point2{math.Abs(eor.Speed.X()), (eor.Speed.Y())}
	}
	if eor.X() > xLimit-eor.W() {
		eor.Delta = floatgeom.Point2{-1 * math.Abs(eor.Speed.X()), (eor.Speed.Y())}
	}
	if eor.Y() < 0 {
		eor.Delta = floatgeom.Point2{eor.Speed.X(), math.Abs(eor.Speed.Y())}
	}
	if eor.Y() > yLimit-eor.H() {
		eor.Delta = floatgeom.Point2{eor.Speed.X(), -1 * math.Abs(eor.Speed.Y())}
	}
	eor.ShiftDelta()
	return 0
}
