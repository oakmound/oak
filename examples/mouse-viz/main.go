package main

import (
	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/debugtools/inputviz"
	"github.com/oakmound/oak/v4/scene"
)

func main() {
	oak.AddScene("mouseviz", scene.Scene{
		Start: func(ctx *scene.Context) {
			bds := ctx.Window.Bounds()
			m := inputviz.Mouse{
				Rect:      floatgeom.NewRect2(0, 0, float64(bds.X()), float64(bds.Y())),
				BaseLayer: -1,
			}
			m.RenderAndListen(ctx, 0)
		},
	})
	oak.Init("mouseviz", func(c oak.Config) (oak.Config, error) {
		c.Screen.Width = 100
		c.Screen.Height = 140
		return c, nil
	})
}
