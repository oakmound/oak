package btn

import (
	"image/color"
	"strconv"

	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
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
	R1           render.Modifiable
	R2           render.Modifiable
	RS           []render.Modifiable
	Cid          event.CID
	Font         *render.Font
	Layer        int
	Text         string
	Children     []Generator
	// This should be a map
	Binding    event.Bindable
	Trigger    string
	Toggle     *bool
	ListChoice *int
	Group      *Group
}

var (
	// A number of these fields could be removed, because they are the zero
	// value, but are left for documentation
	defaultGenerator = Generator{
		X:     0,
		Y:     0,
		W:     1,
		H:     1,
		TxtX:  0,
		TxtY:  0,
		Color: color.RGBA{255, 0, 0, 255},
		Mod:   nil,
		R1:    nil,
		R2:    nil,

		Children: []Generator{},
		Cid:      0,
		Font:     nil,
		Layer:    0,
		Text:     "Button",
		Binding:  nil,
		Trigger:  "MouseClickOn",

		Toggle: nil,
	}
)

// Generate creates a Button from a generator.
func (g Generator) Generate() Btn {
	return g.generate(nil)
}
func (g Generator) generate(parent *Generator) Btn {
	var box render.Modifiable
	if g.Toggle != nil {
		//Handles checks and other toggle situations
		start := "on"
		if !(*g.Toggle) {
			start = "off"
		}
		box = render.NewSwitch(start, map[string]render.Modifiable{
			"on":  g.R1,
			"off": g.R2,
		})
		g.Binding = func(id int, nothing interface{}) int {
			btn := event.GetEntity(id).(Btn)
			if btn.GetRenderable().(*render.Switch).Get() == "on" {
				if g.Group != nil && g.Group.active == btn {
					g.Group.active = nil
				}
				btn.GetRenderable().(*render.Switch).Set("off")
			} else {
				// We can pull this out to seperate binding if group != nil
				if g.Group != nil {
					g.Group.active = btn
					for _, b := range g.Group.members {
						if b.GetRenderable().(*render.Switch).Get() == "on" {
							b.Trigger("MouseClickOn", nil)
						}
					}
				}
				btn.GetRenderable().(*render.Switch).Set("on")

			}
			*g.Toggle = !*g.Toggle

			return 0
		}
		g.Trigger = "MouseClickOn"
	} else if g.ListChoice != nil {

		start := "list" + strconv.Itoa(*g.ListChoice)
		mp := make(map[string]render.Modifiable)
		for i, r := range g.RS {
			mp["list"+strconv.Itoa(i)] = r
		}
		box = render.NewSwitch(start, mp)

		g.Binding = func(id int, button interface{}) int {
			btn := event.GetEntity(id).(*TextBox)
			i := *g.ListChoice
			mEvent := button.(mouse.Event)

			if mEvent.Button == "LeftMouse" {
				i++
				if i == len(g.RS) {
					i = 0
				}

			} else if mEvent.Button == "RightMouse" {
				i--
				if i < 0 {
					i += len(g.RS)
				}
			}

			btn.R.(*render.Switch).Set("list" + strconv.Itoa(i))

			*g.ListChoice = i

			return 0
		}
		g.Trigger = "MouseClickOn"
	} else if g.ProgressFunc != nil {
		box = render.NewGradientBox(int(g.W), int(g.H), g.Color, g.Color2, g.ProgressFunc)
	} else {
		box = render.NewColorBox(int(g.W), int(g.H), g.Color)
	}

	if g.Mod != nil {
		box.Modify(g.Mod)
	}
	font := g.Font
	if font == nil {
		font = render.DefFont()
	}
	// Todo: if no string is defined, don't do this
	btn := NewTextBox(g.Cid, g.X, g.Y, g.W, g.H, g.TxtX, g.TxtY, font, box, g.Layer)
	btn.SetString(g.Text)

	if g.Binding != nil {
		btn.Bind(g.Binding, g.Trigger)
	}

	if g.Group != nil {
		g.Group.members = append(g.Group.members, btn)
	}

	return btn
}

// An Option is used to populate generator fields prior to generation of a button
type Option func(Generator) Generator

// New creates a button with the given options and defaults for all variables not set.
func New(opts ...Option) Btn {
	g := defaultGenerator
	for _, opt := range opts {
		g = opt(g)
	}
	return g.Generate()
}
