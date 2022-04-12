package main

import (
	"image/color"
	"math"
	"math/rand"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

var (
	score1 = 0
	score2 = 0
)

const (
	hitPaddle collision.Label = 1
)

func main() {
	oak.AddScene("pong",
		scene.Scene{Start: func(ctx *scene.Context) {
			newPaddle(ctx, 20, 200, 1)
			newPaddle(ctx, 600, 200, 2)
			newBall(ctx, 320, 240)
			ctx.Draw(render.NewIntText(&score2, 200, 20), 3)
			ctx.Draw(render.NewIntText(&score1, 400, 20), 3)
		}})
	oak.Init("pong")
}

func newBallDelta() floatgeom.Point2 {
	d := floatgeom.Point2{(rand.Float64() - 0.5) * 4, (rand.Float64() - 0.5) * 16}
	if math.Abs(d.X()) < 0.5 {
		d[0] *= 5
	}
	return d
}

func newBall(ctx *scene.Context, x, y float64) {
	ball := entities.New(ctx,
		entities.WithRect(floatgeom.NewRect2WH(x, y, 10, 10)),
		entities.WithColor(color.RGBA{255, 255, 255, 255}),
		entities.WithDrawLayers([]int{2}),
	)
	ball.Delta = newBallDelta()
	event.Bind(ctx, event.Enter, ball, func(ball *entities.Entity, _ event.EnterPayload) event.Response {
		ball.ShiftDelta()
		if collision.HitLabel(ball.Space, hitPaddle) != nil {
			ball.Delta[0] *= -1.1
			ball.Delta[1] += (rand.Float64() - 0.5) * 8
		}
		if ball.X() < 0 || ball.X() > 640 {
			if ball.X() < 0 {
				score1++
			} else {
				score2++
			}
			ball.Delta = newBallDelta()
			ball.SetPos(floatgeom.Point2{320, 240})
		} else if ball.Y() < 0 || ball.Y() > 480-ball.H() {
			ball.Delta[1] = -1 * ball.Delta.Y()
		}
		return 0
	})
}

func newPaddle(ctx *scene.Context, x, y float64, player int) {
	paddle := entities.New(ctx,
		entities.WithRect(floatgeom.NewRect2WH(x, y, 20, 100)),
		entities.WithColor(color.RGBA{255, 255, 255, 255}),
		entities.WithDrawLayers([]int{1}),
		entities.WithLabel(hitPaddle),
	)
	if player == 2 {
		event.Bind(ctx, event.Enter, paddle, enterPaddle(key.UpArrow, key.DownArrow))
	} else {
		event.Bind(ctx, event.Enter, paddle, enterPaddle(key.W, key.S))
	}
}

func enterPaddle(up, down key.Code) func(*entities.Entity, event.EnterPayload) event.Response {
	return func(p *entities.Entity, _ event.EnterPayload) event.Response {
		if oak.IsDown(up) {
			if p.Y() > 0 {
				p.ShiftY(-8)
			}
		} else if oak.IsDown(down) && p.Y() < (480-p.H()) {
			p.ShiftY(8)
		}
		return 0
	}
}
