package particle

import (
	"image/color"
)

type Colorable interface {
	SetStartColor(color.Color, color.Color)
	SetEndColor(color.Color, color.Color)
}

func Color(sc, scr, ec, ecr color.Color) func(Generator) {
	return func(g Generator) {
		c := g.(Colorable)
		c.SetStartColor(sc, scr)
		c.SetEndColor(ec, ecr)
	}
}

type Colorable2 interface {
	SetStartColor2(color.Color, color.Color)
	SetEndColor2(color.Color, color.Color)
}

func Color2(sc, scr, ec, ecr color.Color) func(Generator) {
	return func(g Generator) {
		c := g.(Colorable2)
		c.SetStartColor2(sc, scr)
		c.SetEndColor2(ec, ecr)
	}
}
