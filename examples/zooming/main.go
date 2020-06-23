package main

import (
	"image/color"
	"image/draw"
	"path/filepath"

	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/key"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
)

var (
	zoomOutFactorX = 1.0
	zoomOutFactorY = 1.0
)

func main() {
	oak.Add("demo", func(string, interface{}) {
		s, err := render.LoadSprite("assets", filepath.Join("raw", "mona-lisa.jpg"))
		dlog.ErrorCheck(err)
		zoomer := &zoomR{
			Renderable: s,
			SetFn: func(buff draw.Image, x, y int, c color.Color) {
				x = int(float64(x) / zoomOutFactorX)
				y = int(float64(y) / zoomOutFactorY)
				buff.Set(x, y, c)
			},
		}
		render.Draw(zoomer)
		event.GlobalBind(func(i int, _ interface{}) int {
			if oak.IsDown(key.UpArrow) {
				zoomOutFactorY -= .10
			}
			if oak.IsDown(key.DownArrow) {
				zoomOutFactorY += .10
			}
			if oak.IsDown(key.RightArrow) {
				zoomOutFactorX += .10
			}
			if oak.IsDown(key.LeftArrow) {
				zoomOutFactorX -= .10
			}
			return 0
		}, event.Enter)

	}, func() bool {
		return true
	}, scene.GoTo("demo"))
	oak.Init("demo")
}

type zoomR struct {
	render.Renderable
	SetFn func(buff draw.Image, x, y int, c color.Color)
}

func (z *zoomR) Draw(buff draw.Image) {
	z.DrawOffset(buff, 0, 0)
}
func (z *zoomR) DrawOffset(buff draw.Image, xOff, yOff float64) {
	img := &customImage{buff, z.SetFn}
	z.Renderable.DrawOffset(img, xOff, yOff)
}

type customImage struct {
	draw.Image
	SetFn func(buff draw.Image, x, y int, c color.Color)
}

func (c *customImage) Set(x, y int, col color.Color) {
	c.SetFn(c.Image, x, y, col)
}
