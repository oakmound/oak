package btn

import (
	"fmt"

	"github.com/oakmound/oak/v3/render"
)

//Text sets the text of the button to be generated
func Text(s string) Option {
	return func(g Generator) Generator {
		g.TextPtr = nil
		g.TextStringer = nil
		g.Text = s
		return g
	}
}

// TextPtr sets the text of the button to be generated
// to a string pointer.
func TextPtr(s *string) Option {
	return func(g Generator) Generator {
		g.Text = ""
		g.TextStringer = nil
		g.TextPtr = s
		return g
	}
}

// TextStringer sets the text of the generated button to
// use a fmt.Stringer String call
func TextStringer(s fmt.Stringer) Option {
	return func(g Generator) Generator {
		g.Text = ""
		g.TextPtr = nil
		g.TextStringer = s
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
		if g.Font != nil && (g.Text != "" || g.TextPtr != nil) {
			measure := g.Text
			if g.TextPtr != nil {
				measure = *g.TextPtr
			}
			w := g.Font.MeasureString(measure)
			wf := float64(w.Ceil() + buffer)
			if g.W < wf {
				g.W = wf
			}
		}
		return g
	}
}
