package particle

import (
	"image/color"
)

// A Colorable can have colors set on it.
type Colorable interface {
	SetStartColor(color.Color, color.Color)
	SetEndColor(color.Color, color.Color)
}

// Color sets colors on a Colorable
func Color(sc, scr, ec, ecr color.Color) func(Generator) {
	return func(g Generator) {
		c := g.(Colorable)
		c.SetStartColor(sc, scr)
		c.SetEndColor(ec, ecr)
	}
}

// A Colorable2 can have more colors set on it
type Colorable2 interface {
	SetStartColor2(color.Color, color.Color)
	SetEndColor2(color.Color, color.Color)
}

// Color2 sets more colors on a Colorable2
func Color2(sc, scr, ec, ecr color.Color) func(Generator) {
	return func(g Generator) {
		c := g.(Colorable2)
		c.SetStartColor2(sc, scr)
		c.SetEndColor2(ec, ecr)
	}
}
