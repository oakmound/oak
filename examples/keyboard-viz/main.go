package main

import (
	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/debugtools/inputviz"
	"github.com/oakmound/oak/v3/scene"
)

func main() {
	oak.AddScene("keyviz", scene.Scene{
		Start: func(ctx *scene.Context) {
			m := inputviz.Keyboard{
				Rect:      floatgeom.NewRect2(0, 0, float64(ctx.Window.Width()), float64(ctx.Window.Height())),
				BaseLayer: -1,
			}
			m.RenderAndListen(ctx, 0)
		},
	})
	oak.Init("keyviz", func(c oak.Config) (oak.Config, error) {
		c.Screen.Width = 800
		c.Screen.Height = 300
		return c, nil
	})
}
