package main

import (
	"image/color"
	"math/rand"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

var (
	score1 = 0
	score2 = 0
)

const (
	paddle collision.Label = 1
)

func main() {
	oak.AddScene("pong",
		func(prevScene string, data interface{}) {
			newPaddle(20, 200, 1)
			newPaddle(590, 200, 2)
			newBall(320, 240)
			render.Draw(render.DefFont().NewIntText(&score2, 200, 20), 3)
			render.Draw(render.DefFont().NewIntText(&score1, 400, 20), 3)
		}, func() bool { return true },
		func() (string, *oak.SceneResult) { return "pong", nil })
	oak.Init("pong")
}

func newBall(x, y float64) {
	b := entities.NewMoving(x, y, 10, 10, render.NewColorBox(10, 10, color.RGBA{0, 255, 0, 255}), 0, 0)
	render.Draw(b.R, 2)
	b.Bind(func(id int, nothing interface{}) int {
		if b.Delta.X() == 0 && b.Delta.Y() == 0 {
			b.Delta.SetY((rand.Float64() - 0.5) * 4)
			b.Delta.SetX((rand.Float64() - 0.5) * 16)
			if b.Delta.X() == 0 {
				b.Delta.SetX(8)
			}
		}
		b.ShiftPos(b.Delta.X(), b.Delta.Y())
		if collision.HitLabel(b.Space, paddle) != nil {
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
		} else if b.Y() < 0 || b.Y() > 440-b.H {
			b.Delta.SetY(-1 * b.Delta.Y())
		}
		return 0
	}, "EnterFrame")
}

func newPaddle(x, y float64, player int) {
	p := entities.NewMoving(x, y, 20, 100, render.NewColorBox(20, 100, color.RGBA{255, 0, 0, 255}), 0, 0)
	p.Speed.SetY(4)
	render.Draw(p.R, 1)
	p.Space.UpdateLabel(paddle)
	if player == 1 {
		p.Bind(enterPaddle("UpArrow", "DownArrow"), "EnterFrame")
	} else {
		p.Bind(enterPaddle("W", "S"), "EnterFrame")
	}
	p.SetPos(x, y)
}

func enterPaddle(up, down string) func(int, interface{}) int {
	return func(id int, nothing interface{}) int {
		p := event.GetEntity(id).(*entities.Moving)
		p.Delta.SetY(0)
		if oak.IsDown(up) {
			p.Delta.SetY(-p.Speed.Y())
		} else if oak.IsDown(down) {
			p.Delta.SetY(p.Speed.Y())
		}
		p.ShiftY(p.Delta.Y())
		if p.Y() < 0 || p.Y() > (440-p.H) {
			p.ShiftY(-p.Delta.Y())
		}
		return 0
	}
}
