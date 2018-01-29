package main

import (
	"image/color"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
)

const (
	Enemy  collision.Label = 1
	Player collision.Label = 2
)

var (
	playerAlive = true
)

func main() {
	oak.Add("tds", func(string, interface{}) {
		playerAlive = true
		char := entities.NewMoving(100, 100, 32, 32,
			render.NewColorBox(32, 32, color.RGBA{0, 255, 0, 255}),
			nil, 0, 0)

		char.Speed = physics.NewVector(5, 5)
		render.Draw(char.R)

		char.Bind(func(id int, _ interface{}) int {
			char := event.GetEntity(id).(*entities.Moving)
			char.Delta.Zero()
			if oak.IsDown(key.W) {
				char.Delta.ShiftY(-char.Speed.Y())
			}
			if oak.IsDown(key.A) {
				char.Delta.ShiftX(-char.Speed.X())
			}
			if oak.IsDown(key.S) {
				char.Delta.ShiftY(char.Speed.Y())
			}
			if oak.IsDown(key.D) {
				char.Delta.ShiftX(char.Speed.X())
			}
			char.ShiftPos(char.Delta.X(), char.Delta.Y())
			hit := char.HitLabel(Enemy)
			if hit != nil {
				playerAlive = false
			}

			return 0
		}, event.Enter)

	}, func() bool {
		return playerAlive
	}, func() (string, *scene.Result) {
		return "tds", nil
	})
	oak.Init("tds")
}
