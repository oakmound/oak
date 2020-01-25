package btn

import "github.com/oakmound/oak/v2/render"

//Text sets the text of the button to be generated
func Text(s string) Option {
	return func(g Generator) Generator {
		g.Text = s
		return g
	}
}

//Font sets the font for the text of the button to be generated
func Font(f *render.Font) Option {
	return func(g Generator) Generator {
		g.Font = f
		return g
	}
}

//TxtOff sets the text offset  of the button generator from the bottom left
func TxtOff(x, y float64) Option {
	return func(g Generator) Generator {
		g.TxtX = x
		g.TxtY = y
		return g
	}
}
