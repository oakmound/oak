package main

import (
	"image/color"
	"math"
	"math/rand"

	oak "github.com/oakmound/oak/v3"
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
			render.Draw(render.DefaultFont().NewIntText(&score2, 200, 20), 3)
			render.Draw(render.DefaultFont().NewIntText(&score1, 400, 20), 3)
		}})
	oak.Init("pong", func(c oak.Config) (oak.Config, error) {
		c.DrawFrameRate = 120
		return c, nil
	})
}

func newBall(ctx *scene.Context, x, y float64) {
	b := entities.NewMoving(x, y, 10, 10, render.NewColorBoxR(10, 10, color.RGBA{255, 255, 255, 255}), nil, 0, 0)
	render.Draw(b.R, 2)
	event.GlobalBind(ctx, event.Enter, func(_ event.EnterPayload) event.Response {
		if b.Delta.X() == 0 && b.Delta.Y() == 0 {
			b.Delta.SetY((rand.Float64() - 0.5) * 4)
			b.Delta.SetX((rand.Float64() - 0.5) * 16)
			if math.Abs(b.Delta.X()) < 0.1 {
				b.Delta.SetX(8)
			}
		}
		b.ShiftPos(b.Delta.X(), b.Delta.Y())
		if collision.HitLabel(b.Space, hitPaddle) != nil {
			b.Delta.SetX(-1.1 * b.Delta.X())
			b.Delta.SetY(b.Delta.Y() + (rand.Float64()-0.5)*8)
		}
		if b.X() < 0 || b.X() > 640 {
			if b.X() < 0 {
				score1++
			} else {
				score2++
			}
			b.Delta.SetX(0)
			b.Delta.SetY(0)
			b.SetPos(320, 240)
		} else if b.Y() < 0 || b.Y() > 480-b.H {
			b.Delta.SetY(-1 * b.Delta.Y())
		}
		return 0
	})
}

func newPaddle(ctx *scene.Context, x, y float64, player int) {
	p := entities.NewMoving(x, y, 20, 100, render.NewColorBoxR(20, 100, color.RGBA{255, 255, 255, 255}), nil, 0, 0)
	p.Speed.SetY(8)
	render.Draw(p.R, 1)
	p.Space.UpdateLabel(hitPaddle)
	if player == 1 {
		event.Bind(ctx, event.Enter, p, enterPaddle(key.UpArrowStr, key.DownArrowStr))
	} else {
		event.Bind(ctx, event.Enter, p, enterPaddle(key.WStr, key.SStr))
	}
}

func enterPaddle(up, down string) func(*entities.Moving, event.EnterPayload) event.Response {
	return func(p *entities.Moving, _ event.EnterPayload) event.Response {
		p.Delta.SetY(0)
		if oak.IsDown(up) {
			p.Delta.SetY(-p.Speed.Y())
		} else if oak.IsDown(down) {
			p.Delta.SetY(p.Speed.Y())
		}
		p.ShiftY(p.Delta.Y())
		if p.Y() < 0 || p.Y() > (480-p.H) {
			p.ShiftY(-p.Delta.Y())
		}
		return 0
	}
}
