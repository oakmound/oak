package main

import (
	"embed"
	"image/color"
	"image/draw"
	"path/filepath"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

var (
	zoomOutFactorX = 1.0
	zoomOutFactorY = 1.0
)

func main() {
	oak.AddScene("demo", scene.Scene{Start: func(*scene.Context) {
		render.Draw(render.NewText("Controls: Arrow keys", 500, 440))

		// Get an image that we will illustrate zooming with later
		s, err := render.LoadSprite(filepath.Join("assets", "mona-lisa.jpg"))
		dlog.ErrorCheck(err)

		// See the zoomR definition lower, wrap your renderable with a definition of how to zoom.
		zoomer := &zoomR{
			Renderable: s,
			SetFn: func(buff draw.Image, x, y int, c color.Color) {
				x = int(float64(x) * zoomOutFactorX)
				y = int(float64(y) * zoomOutFactorY)
				buff.Set(x, y, c)
			},
		}
		render.Draw(zoomer)

		// To illustrate zooming allow for arrow keys to control the main zoomable renderable.
		event.GlobalBind(event.Enter, func(i event.CID, _ interface{}) int {
			if oak.IsDown(key.UpArrow) {
				zoomOutFactorY *= .98
			}
			if oak.IsDown(key.DownArrow) {
				zoomOutFactorY *= 1.02
			}
			if oak.IsDown(key.RightArrow) {
				zoomOutFactorX *= 1.02
			}
			if oak.IsDown(key.LeftArrow) {
				zoomOutFactorX *= .98
			}

			return 0
		})

	}})
	oak.SetFS(assets)
	oak.Init("demo")
}

//go:embed assets
var assets embed.FS

// zoomR wraps a renderable with a function that details
type zoomR struct {
	render.Renderable
	SetFn func(buff draw.Image, x, y int, c color.Color)
}

// Draw to draw the zoomR by creating a customImage and applying the set funcitonality.
func (z *zoomR) Draw(buff draw.Image, xOff, yOff float64) {
	img := &customImage{buff, z.SetFn}
	z.Renderable.Draw(img, xOff, yOff)
}

type customImage struct {
	draw.Image
	SetFn func(buff draw.Image, x, y int, c color.Color)
}

func (c *customImage) Set(x, y int, col color.Color) {
	c.SetFn(c.Image, x, y, col)
}
