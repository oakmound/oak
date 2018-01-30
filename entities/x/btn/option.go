package btn

import (
	"image/color"

	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
)

// And combines a variadic number of options
func And(opts ...Option) Option {
	return func(g Generator) Generator {
		for _, opt := range opts {
			g = opt(g)
		}
		return g
	}
}

//Width sets the Width of the button to be generated
func Width(w float64) Option {
	return func(g Generator) Generator {
		g.W = w
		return g
	}
}

//Height sets the Height of the button to be generated
func Height(h float64) Option {
	return func(g Generator) Generator {
		g.H = h
		return g
	}
}

//Pos sets the position of the button  to be generated
func Pos(x, y float64) Option {
	return func(g Generator) Generator {
		g.X = x
		g.Y = y
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

//CID sets the starting CID of the button to be generated
func CID(c event.CID) Option {
	return func(g Generator) Generator {
		g.Cid = c
		return g
	}
}

//Color sets the colorboxes color for the button to be generated
func Color(c color.Color) Option {
	return func(g Generator) Generator {
		g.Color = c
		return g
	}
}

// VGradient creates a vertical color gradient for the btn
func VGradient(c1, c2 color.Color) Option {
	return func(g Generator) Generator {
		g.Color = c1
		g.Color2 = c2
		g.ProgressFunc = render.VerticalProgress
		return g
	}
}

//Mod sets the modifications to apply to the initial color box for the button to be generated
func Mod(m mod.Transform) Option {
	return func(g Generator) Generator {
		g.Mod = m
		return g
	}
}

// AndMod combines the input modification with whatever existing modifications
// exist for the generator, as opposed to Mod which resets previous modifications.
func AndMod(m mod.Transform) Option {
	return func(g Generator) Generator {
		if g.Mod == nil {
			g.Mod = m
		} else {
			g.Mod = mod.And(g.Mod, m)
		}
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

//Layer sets the layer of the button to be generated
func Layer(l int) Option {
	return func(g Generator) Generator {
		g.Layer = l
		return g
	}
}

//Text sets the text of the button to be generated
func Text(s string) Option {
	return func(g Generator) Generator {
		g.Text = s
		return g
	}
}

// Toggle sets that the type of the button toggles between two
// modifiables when it is clicked. The boolean behind isChecked
// is updated according to the state of the button.
func Toggle(r1, r2 render.Modifiable, isChecked *bool) Option {
	return func(g Generator) Generator {
		g.R1 = r1.Copy()
		g.R2 = r2.Copy()
		g.Toggle = isChecked
		return g
	}
}

// ToggleGroup sets the group that this button is linked with
func ToggleGroup(gr *Group) Option {
	return func(g Generator) Generator {
		g.Group = gr
		return g
	}
}

// ToggleList sets the togglable choices for a button
func ToggleList(chosen *int, rs ...render.Modifiable) Option {
	return func(g Generator) Generator {
		g.ListChoice = chosen
		g.RS = rs
		return g
	}
}

//Binding sets the Binding of the button to be generated
func Binding(bnd event.Bindable) Option {
	return func(g Generator) Generator {
		g.Binding = bnd
		return g
	}
}

//Trigger sets the trigger for the Binding on the button to be generated
func Trigger(s string) Option {
	return func(g Generator) Generator {
		g.Trigger = s
		return g
	}
}
