package btn

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/render/mod"
	"github.com/oakmound/oak/v3/shape"
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
	Bindings     []func(caller Btn) event.Binding
	Trigger      string
	Toggle       *bool
	ListChoice   *int
	Group        *Group
	AllowRevert  bool
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

		Children: []Generator{},
		Cid:      0,
		Font:     nil,
		Layers:   []int{0},
		Text:     "Button",
		Trigger:  "MouseClickOn",

		Toggle: nil,
	}
}

// Generate creates a Button from a generator.
func (g Generator) Generate() Btn {
	return g.generate(nil)
}

func (g Generator) generate(parent *Generator) Btn {
	var box render.Modifiable
	// handle different renderable options that could be passed to the generator
	switch {
	case g.Toggle != nil:
		//Handles checks and other toggle situations
		start := "on"
		if !(*g.Toggle) {
			start = "off"
		}
		if _, ok := g.R1.(*render.Reverting); !ok {
			g.R1 = render.NewReverting(g.R1)
		}
		if _, ok := g.R2.(*render.Reverting); !ok {
			g.R2 = render.NewReverting(g.R2)
		}
		box = render.NewSwitch(start, map[string]render.Modifiable{
			"on":  g.R1,
			"off": g.R2,
		})
		g.Bindings = append(g.Bindings, func(caller Btn) event.Binding {
			// TODO not default
			return event.Bind(event.DefaultBus, mouse.ClickOn, caller, toggleFxn(g))
		})
	case g.ListChoice != nil:

		start := "list" + strconv.Itoa(*g.ListChoice)
		mp := make(map[string]render.Modifiable)
		for i, r := range g.RS {
			if _, ok := r.(*render.Reverting); !ok {
				r = render.NewReverting(r)
			}
			mp["list"+strconv.Itoa(i)] = r
		}

		box = render.NewSwitch(start, mp)

		g.Bindings = append(g.Bindings, func(caller Btn) event.Binding {
			// TODO not default
			return event.Bind(event.DefaultBus, mouse.ClickOn, caller, listFxn(g))
		})
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

	if g.AllowRevert {
		box = render.NewReverting(box)
	}

	if g.Mod != nil {
		box.Modify(g.Mod)
	}
	font := g.Font
	if font == nil {
		font = render.DefaultFont()
	}
	var btn Btn
	if g.Text != "" {
		txtbx := NewTextBox(g.Cid, g.X, g.Y, g.W, g.H, g.TxtX, g.TxtY, font, box, g.Layers...)
		txtbx.SetString(g.Text)
		txtbx.Space.Label = g.Label
		btn = txtbx
	} else if g.TextPtr != nil {
		txtbx := NewTextBox(g.Cid, g.X, g.Y, g.W, g.H, g.TxtX, g.TxtY, font, box, g.Layers...)
		txtbx.SetStringPtr(g.TextPtr)
		txtbx.Space.Label = g.Label
		btn = txtbx
	} else if g.TextStringer != nil {
		txtbx := NewTextBox(g.Cid, g.X, g.Y, g.W, g.H, g.TxtX, g.TxtY, font, box, g.Layers...)
		txtbx.SetStringer(g.TextStringer)
		txtbx.Space.Label = g.Label
		btn = txtbx
	} else {
		bx := NewBox(g.Cid, g.X, g.Y, g.W, g.H, box, g.Layers...)
		bx.Space.Label = g.Label
		btn = bx
	}

	// TODO: this is impossible with how we've done generics

	// Update underlying mousecollision binding to only respect clicks in the shape.
	// If a finer control is needed then it may make sense to use this as a starting off point
	// instead of expanding this section.
	// if g.Shape != nil {
	// 	// extract keys prior to loop as the map will be permuted by the following operations
	// 	keys := make([]string, 0, len(g.Bindings))
	// 	for k := range g.Bindings {
	// 		// We only really care about mouse events.
	// 		// In some ways this is dangerous of an implementer has defined events that start with mouse...
	// 		// but in that case they might not use g.Shape anyways.
	// 		if !strings.HasPrefix(k, "Mouse") {
	// 			continue
	// 		}
	// 		keys = append(keys, k)
	// 	}
	// 	for _, k := range keys {
	// 		curBind := g.Bindings[k]
	// 		if curBind == nil {
	// 			continue
	// 		}
	// 		// This could cause issues with name collisions but its unlikely and documentation should help make it even more unlikely.
	// 		filteredK := "Filtered" + k
	// 		g.Bindings[filteredK] = g.Bindings[k]
	// 		g.Bindings[k] = []event.Bindable{
	// 			func(id event.CallerID, button interface{}) int {
	// 				btn := id.E().(Btn)
	// 				mEvent, ok := button.(*mouse.Event)
	// 				// If the passed event is not a mouse event dont filter on location.
	// 				// Main current use case is for nil events passed via simulated clicks.
	// 				if !ok {
	// 					btn.Trigger(filteredK, button)
	// 				}
	// 				bSpace := btn.GetSpace().Bounds()
	// 				if g.Shape.In(int(mEvent.X()-bSpace.Min.X()), int(mEvent.Y()-bSpace.Min.Y()), int(bSpace.W()), int(bSpace.H())) {
	// 					btn.Trigger(filteredK, mEvent)
	// 				}
	// 				return 0
	// 			},
	// 		}
	// 	}
	// }

	for _, binding := range g.Bindings {
		binding(btn)
	}

	// TODO: not default
	err := mouse.PhaseCollision(btn.GetSpace(), event.DefaultBus)
	dlog.ErrorCheck(err)

	if g.Group != nil {
		g.Group.members = append(g.Group.members, btn)
	}

	return btn
}

// An Option is used to populate generator fields prior to generation of a button
type Option func(Generator) Generator

// New creates a button with the given options and defaults for all variables not set.
func New(opts ...Option) Btn {
	g := defGenerator()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		g = opt(g)
	}
	return g.Generate()
}

type switcher interface {
	Get() string
	Set(string) error
}

// toggleFxn sets up the mouseclick binding for toggle buttons created for goreport cyclo decrease
func toggleFxn(g Generator) func(btn Btn, payload *mouse.Event) event.Response {
	return func(btn Btn, payload *mouse.Event) event.Response {
		if btn.GetRenderable().(switcher).Get() == "on" {
			if g.Group != nil && g.Group.active == btn {
				g.Group.active = nil
			}
			btn.GetRenderable().(switcher).Set("off")
		} else {
			// We can pull this out to separate binding if group != nil
			if g.Group != nil {
				g.Group.active = btn
				for _, b := range g.Group.members {
					if b.GetRenderable().(switcher).Get() == "on" {
						event.TriggerForCallerOn(event.DefaultBus, b.CID(), mouse.ClickOn, payload)
					}
				}
			}
			btn.GetRenderable().(switcher).Set("on")

		}
		*g.Toggle = !*g.Toggle

		return 0
	}
}

// listFxn sets up the mouseclick binding for list buttons created for goreport cyclo reduction
func listFxn(g Generator) func(btn Btn, payload *mouse.Event) event.Response {
	return func(btn Btn, payload *mouse.Event) event.Response {
		i := *g.ListChoice
		if payload.Button == mouse.ButtonLeft {
			i++
			if i == len(g.RS) {
				i = 0
			}

		} else if payload.Button == mouse.ButtonRight {
			i--
			if i < 0 {
				i += len(g.RS)
			}
		}

		btn.GetRenderable().(*render.Switch).Set("list" + strconv.Itoa(i))

		*g.ListChoice = i

		return 0
	}
}
