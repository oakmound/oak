package btn

import (
	"image/color"

	"github.com/oakmound/oak/v4/collision"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/mouse"
	"github.com/oakmound/oak/v4/scene"

	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/render/mod"
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

// Clear resets the button to be empty
func Clear() Option {
	return func(g Generator) Generator {
		return Generator{}
	}
}

// Width sets the Width of the button to be generated
func Width(w float64) Option {
	return func(g Generator) Generator {
		g.W = w
		return g
	}
}

// Height sets the Height of the button to be generated
func Height(h float64) Option {
	return func(g Generator) Generator {
		g.H = h
		return g
	}
}

// Pos sets the position of the button  to be generated
func Pos(x, y float64) Option {
	return func(g Generator) Generator {
		g.X = x
		g.Y = y
		return g
	}
}

// Offset increments the position of the button to be generated
func Offset(x, y float64) Option {
	return func(g Generator) Generator {
		g.X += x
		g.Y += y
		return g
	}
}

// CID sets the starting CID of the button to be generated
func CID(c event.CallerID) Option {
	return func(g Generator) Generator {
		g.Cid = c
		return g
	}
}

// Color sets the colorboxes color for the button to be generated
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

// Mod sets the modifications to apply to the initial color box for the button to be generated
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

// Layers sets the layer of the button to be generated
func Layers(ls ...int) Option {
	return func(g Generator) Generator {
		g.Layers = ls
		return g
	}
}

// Renderable sets a renderable to use as a base image for the button.
// Not compatible with Color / Toggle.
func Renderable(r render.Modifiable) Option {
	return func(g Generator) Generator {
		g.R = r
		return g
	}
}

// Binding appends a function to be called when a specific event
// is triggered.
func Binding[Payload any](ev event.EventID[Payload], bnd event.Bindable[*entities.Entity, Payload]) Option {
	return func(g Generator) Generator {
		g.Bindings = append(g.Bindings, func(ctx *scene.Context, caller *entities.Entity) event.Binding {
			// TODO: not default
			return event.Bind(ctx, ev, caller, bnd)
		})
		return g
	}
}

// Click appends a function to be called when the button is clicked on.
func Click(bnd event.Bindable[*entities.Entity, *mouse.Event]) Option {
	return Binding(mouse.ClickOn, bnd)
}

func Label(l collision.Label) Option {
	return func(g Generator) Generator {
		g.Label = l
		return g
	}
}
