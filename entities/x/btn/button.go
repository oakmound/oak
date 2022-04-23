package btn

import (
	"fmt"
	"image/color"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/collision"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/mouse"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/render/mod"
	"github.com/oakmound/oak/v4/scene"
	"github.com/oakmound/oak/v4/shape"
)

// A Generator defines the variables used to create buttons from optional arguments
type Generator struct {
	X, Y         float64
	W, H         float64
	TxtX, TxtY   float64
	Color        color.Color
	Color2       color.Color
	ProgressFunc func(x, y, w, h int) float64
	Mod          mod.Transform
	R            render.Modifiable
	R1           render.Modifiable
	R2           render.Modifiable
	RS           []render.Modifiable
	Cid          event.CallerID
	Font         *render.Font
	Layers       []int
	Text         string
	TextPtr      *string
	TextStringer fmt.Stringer
	Children     []Generator
	Bindings     []func(ctx *scene.Context, caller *entities.Entity) event.Binding
	Trigger      string
	Shape        shape.Shape
	Label        collision.Label
}

func defGenerator() Generator {
	// A number of these fields could be removed, because they are the zero
	// value, but are left for documentation
	return Generator{
		X:     0,
		Y:     0,
		W:     1,
		H:     1,
		TxtX:  0,
		TxtY:  0,
		Color: color.RGBA{255, 0, 0, 255},
		Mod:   nil,
		R:     nil,
		R1:    nil,
		R2:    nil,

		Font:    nil,
		Layers:  []int{0},
		Text:    "",
		Trigger: "MouseClickOn",
	}
}

// Generate creates a Button from a generator.
func (g Generator) Generate(ctx *scene.Context) *entities.Entity {
	var box render.Modifiable
	// handle different renderable options that could be passed to the generator
	switch {
	case g.R != nil:
		box = g.R
	case g.ProgressFunc != nil:
		box = render.NewGradientBox(int(g.W), int(g.H), g.Color, g.Color2, g.ProgressFunc)
		if g.Shape != nil {
			g.Mod = mod.SafeAnd(mod.CutShape(g.Shape), g.Mod)
		}
	default:
		box = render.NewColorBox(int(g.W), int(g.H), g.Color)
		if g.Shape != nil {
			g.Mod = mod.SafeAnd(mod.CutShape(g.Shape), g.Mod)
		}
	}

	entOpts := []entities.Option{
		entities.WithRenderable(box),
		entities.WithMod(g.Mod),
		entities.WithRect(floatgeom.NewRect2WH(g.X, g.Y, g.W, g.H)),
		entities.WithLabel(g.Label),
		entities.WithDrawLayers(g.Layers),
		entities.WithUseMouseTree(true),
	}

	font := g.Font
	if font == nil {
		font = render.DefaultFont()
	}
	childLayers := make([]int, len(g.Layers))
	copy(childLayers, g.Layers)
	if len(childLayers) != 0 {
		childLayers[len(childLayers)-1]++
	}
	if g.Text != "" {
		entOpts = append(entOpts, entities.WithChild(
			entities.WithRenderable(font.NewText(g.Text, g.TxtX, g.TxtY)),
			entities.WithDrawLayers(childLayers),
		))
	} else if g.TextPtr != nil {
		entOpts = append(entOpts, entities.WithChild(
			entities.WithRenderable(font.NewStrPtrText(g.TextPtr, g.TxtX, g.TxtY)),
			entities.WithDrawLayers(childLayers),
		))
	} else if g.TextStringer != nil {
		entOpts = append(entOpts, entities.WithChild(
			entities.WithRenderable(font.NewStringerText(g.TextStringer, g.TxtX, g.TxtY)),
			entities.WithDrawLayers(childLayers),
		))
	}

	btn := entities.New(ctx, entOpts...)

	for _, binding := range g.Bindings {
		binding(ctx, btn)
	}

	mouse.PhaseCollision(btn.Space, ctx.Handler)

	return btn
}

// An Option is used to populate generator fields prior to generation of a button
type Option func(Generator) Generator

// New creates a button with the given options and defaults for all variables not set.
func New(ctx *scene.Context, opts ...Option) *entities.Entity {
	g := defGenerator()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		g = opt(g)
	}
	return g.Generate(ctx)
}
