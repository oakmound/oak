package btn

import (
	"github.com/oakmound/oak/v2/render"
)

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

// FitText adjusts a btn's width, given it has text and font defined, to
// be large enough for the given text plus the provided buffer
func FitText(buffer int) Option {
	return func(g Generator) Generator {
		if g.Font != nil && g.Text != "" {
			w := g.Font.MeasureString(g.Text)
			wf := float64(w.Ceil() + buffer)
			if g.W < wf {
				g.W = wf
			}
		}
		return g
	}
}
