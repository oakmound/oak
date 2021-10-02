package inputviz

import (
	"image/color"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

type KeyboardLayout interface {
	KeyRect(k string) floatgeom.Rect2
}

type LayoutQWERTY struct {
	Bounds floatgeom.Rect2
}

func (l LayoutQWERTY) KeyRect(k string) floatgeom.Rect2 {
	// max row = 5.1, drawn down to 6
	row := 0.0
	// max col = 21.2, drawn right to 22.1
	col := 0.0
	width := 0.9
	height := 0.9
	switch k {
	case key.Escape:
	case key.F1:
		col = 2.0
	case key.F2:
		col = 3.0
	case key.F3:
		col = 4.0
	case key.F4:
		col = 5.0
	case key.F5:
		col = 6.5
	case key.F6:
		col = 7.5
	case key.F7:
		col = 8.5
	case key.F8:
		col = 9.5
	case key.F9:
		col = 11.0
	case key.F10:
		col = 12.0
	case key.F11:
		col = 13.0
	case key.F12:
		col = 14.0
	// key.PrintScreen
	// key.ScrollLock
	case key.Pause:
		col = 17.1
	//case key.Tilde:
	//	row = 1.1
	case key.One:
		row = 1.1
		col = 1.0
	case key.Two:
		row = 1.1
		col = 2.0
	case key.Three:
		row = 1.1
		col = 3.0
	case key.Four:
		row = 1.1
		col = 4.0
	case key.Five:
		row = 1.1
		col = 5.0
	case key.Six:
		row = 1.1
		col = 6.0
	case key.Seven:
		row = 1.1
		col = 7.0
	case key.Eight:
		row = 1.1
		col = 8.0
	case key.Nine:
		row = 1.1
		col = 9.0
	case key.Zero:
		row = 1.1
		col = 10.0
	case key.HyphenMinus:
		row = 1.1
		col = 11.0
	case key.EqualSign:
		row = 1.1
		col = 12.0
	case key.DeleteBackspace:
		row = 1.1
		col = 13.0
		width = 1.9
	case key.Insert:
		row = 1.1
		col = 15.1
	case key.Home:
		row = 1.1
		col = 16.1
	case key.PageUp:
		row = 1.1
		col = 17.1
	case key.KeypadNumLock:
		row = 1.1
		col = 18.2
	case key.KeypadSlash:
		row = 1.1
		col = 19.2
	case key.KeypadAsterisk:
		row = 1.1
		col = 20.2
	case key.KeypadHyphenMinus:
		row = 1.1
		col = 21.2
	case key.Tab:
		row = 2.1
		width = 1.4
	case key.Q:
		row = 2.1
		col = 1.5
	case key.W:
		row = 2.1
		col = 2.5
	case key.E:
		row = 2.1
		col = 3.5
	case key.R:
		row = 2.1
		col = 4.5
	case key.T:
		row = 2.1
		col = 5.5
	case key.Y:
		row = 2.1
		col = 6.5
	case key.U:
		row = 2.1
		col = 7.5
	case key.I:
		row = 2.1
		col = 8.5
	case key.O:
		row = 2.1
		col = 9.5
	case key.P:
		row = 2.1
		col = 10.5
	case key.LeftSquareBracket:
		row = 2.1
		col = 11.5
	case key.RightSquareBracket:
		row = 2.1
		col = 12.5
	case key.Backslash:
		row = 2.1
		col = 13.5
		width = 1.4
	case key.DeleteForward:
		row = 2.1
		col = 15.1
	case key.End:
		row = 2.1
		col = 16.1
	case key.PageDown:
		row = 2.1
		col = 17.1
	case key.Keypad7:
		row = 2.1
		col = 18.2
	case key.Keypad8:
		row = 2.1
		col = 19.2
	case key.Keypad9:
		row = 2.1
		col = 20.2
	case key.KeypadPlusSign:
		row = 2.1
		col = 21.2
		height = 1.9
	case key.CapsLock:
		row = 3.1
		width = 1.4
	case key.A:
		row = 3.1
		col = 1.5
	case key.S:
		row = 3.1
		col = 2.5
	case key.D:
		row = 3.1
		col = 3.5
	case key.F:
		row = 3.1
		col = 4.5
	case key.G:
		row = 3.1
		col = 5.5
	case key.H:
		row = 3.1
		col = 6.5
	case key.J:
		row = 3.1
		col = 7.5
	case key.K:
		row = 3.1
		col = 8.5
	case key.L:
		row = 3.1
		col = 9.5
	case key.Semicolon:
		row = 3.1
		col = 10.5
	case key.Apostrophe:
		row = 3.1
		col = 11.5
	case key.ReturnEnter:
		row = 3.1
		col = 12.5
		width = 2.4
	case key.Keypad4:
		row = 3.1
		col = 18.2
	case key.Keypad5:
		row = 3.1
		col = 19.2
	case key.Keypad6:
		row = 3.1
		col = 20.2
	case key.LeftShift:
		row = 4.1
		width = 1.9
	case key.Z:
		row = 4.1
		col = 2.0
	case key.X:
		row = 4.1
		col = 3.0
	case key.C:
		row = 4.1
		col = 4.0
	case key.V:
		row = 4.1
		col = 5.0
	case key.B:
		row = 4.1
		col = 6.0
	case key.N:
		row = 4.1
		col = 7.0
	case key.M:
		row = 4.1
		col = 8.0
	case key.Comma:
		row = 4.1
		col = 9.0
	case key.FullStop:
		row = 4.1
		col = 10.0
	case key.Slash:
		row = 4.1
		col = 11.0
	case key.RightShift:
		row = 4.1
		col = 12.0
		width = 2.9
	case key.UpArrow:
		row = 4.1
		col = 16.1
	case key.Keypad1:
		row = 4.1
		col = 18.2
	case key.Keypad2:
		row = 4.1
		col = 19.2
	case key.Keypad3:
		row = 4.1
		col = 20.2
	case key.KeypadEnter:
		row = 4.1
		col = 21.2
		height = 1.9
	case key.LeftControl:
		row = 5.1
		width = 1.4
	case key.LeftGUI:
		row = 5.1
		col = 1.5
	case key.LeftAlt:
		row = 5.1
		col = 2.5
		width = 1.4
	case key.Spacebar:
		row = 5.1
		col = 4.0
		width = 6.9
	case key.RightAlt:
		row = 5.1
		col = 11.0
		width = 1.4
	case key.RightGUI:
		row = 5.1
		col = 12.5
	case key.RightControl:
		row = 5.1
		col = 13.5
		width = 1.4
	case key.LeftArrow:
		row = 5.1
		col = 15.1
	case key.DownArrow:
		row = 5.1
		col = 16.1
	case key.RightArrow:
		row = 5.1
		col = 17.1
	case key.Keypad0:
		row = 5.1
		col = 18.2
		width = 1.9
	case key.KeypadPeriod:
		row = 5.1
		col = 20.2
	default:
		return floatgeom.Rect2{}
	}
	w, h := l.Bounds.W(), l.Bounds.H()
	rowHeight := h / 6.0
	colWidth := w / 22.1

	x := col * colWidth
	y := row * rowHeight
	keyHeight := height * rowHeight
	keyWidth := width * colWidth

	return floatgeom.NewRect2WH(x, y, keyWidth, keyHeight)
}

var defaultColors = map[string]color.Color{}

type Keyboard struct {
	Rect      floatgeom.Rect2
	BaseLayer int
	Colors    map[string]color.Color
	KeyboardLayout

	event.CID
	ctx *scene.Context

	rs map[string]*render.Switch
}

func (k *Keyboard) Init() event.CID {
	k.CID = k.ctx.CallerMap.NextID(k)
	return k.CID
}

func (k *Keyboard) RenderAndListen(ctx *scene.Context, layer int) error {
	k.ctx = ctx
	k.Init()

	if k.Rect.W() == 0 || k.Rect.H() == 0 {
		k.Rect.Max = k.Rect.Min.Add(floatgeom.Point2{320, 180})
	}
	if k.KeyboardLayout == nil {
		k.KeyboardLayout = LayoutQWERTY{
			Bounds: k.Rect,
		}
	}
	if k.Colors == nil {
		k.Colors = defaultColors
	}

	k.rs = make(map[string]*render.Switch)

	for kv := range allKeys {
		rect := k.KeyboardLayout.KeyRect(kv)
		if rect == (floatgeom.Rect2{}) {
			continue
		}
		pressedColor := color.RGBA{255, 255, 255, 255}
		var unpressedColor color.Color = color.RGBA{160, 160, 160, 255}
		if c, ok := k.Colors[kv]; ok {
			unpressedColor = c
		}
		r := render.NewSwitch("released", map[string]render.Modifiable{
			"pressed":  render.NewColorBox(int(rect.W()), int(rect.H()), pressedColor),
			"released": render.NewColorBox(int(rect.W()), int(rect.H()), unpressedColor),
		})
		r.SetPos(rect.Min.X(), rect.Min.Y())
		k.rs[kv] = r
		if k.BaseLayer == -1 {
			ctx.DrawStack.Draw(r, layer)
		} else {
			ctx.DrawStack.Draw(r, k.BaseLayer, layer)
		}
	}

	k.Bind(key.Down, key.Binding(func(id event.CID, ev key.Event) int {
		kb, _ := k.ctx.CallerMap.GetEntity(id).(*Keyboard)
		btn := ev.Code.String()[4:]
		if kb.rs[btn] == nil {

			return 0
		}
		kb.rs[btn].Set("pressed")
		return 0
	}))
	k.Bind(key.Up, key.Binding(func(id event.CID, ev key.Event) int {
		kb, _ := k.ctx.CallerMap.GetEntity(id).(*Keyboard)
		btn := ev.Code.String()[4:]
		if kb.rs[btn] == nil {
			return 0
		}
		kb.rs[btn].Set("released")
		return 0
	}))

	return nil
}

func (k *Keyboard) Destroy() {
	k.UnbindAll()
	for _, r := range k.rs {
		r.Undraw()
	}
}

var allKeys = map[string]struct{}{
	key.Unknown: {},

	key.A: {},
	key.B: {},
	key.C: {},
	key.D: {},
	key.E: {},
	key.F: {},
	key.G: {},
	key.H: {},
	key.I: {},
	key.J: {},
	key.K: {},
	key.L: {},
	key.M: {},
	key.N: {},
	key.O: {},
	key.P: {},
	key.Q: {},
	key.R: {},
	key.S: {},
	key.T: {},
	key.U: {},
	key.V: {},
	key.W: {},
	key.X: {},
	key.Y: {},
	key.Z: {},

	key.One:   {},
	key.Two:   {},
	key.Three: {},
	key.Four:  {},
	key.Five:  {},
	key.Six:   {},
	key.Seven: {},
	key.Eight: {},
	key.Nine:  {},
	key.Zero:  {},

	key.ReturnEnter:        {},
	key.Escape:             {},
	key.DeleteBackspace:    {},
	key.Tab:                {},
	key.Spacebar:           {},
	key.HyphenMinus:        {},
	key.EqualSign:          {},
	key.LeftSquareBracket:  {},
	key.RightSquareBracket: {},
	key.Backslash:          {},
	key.Semicolon:          {},
	key.Apostrophe:         {},
	key.GraveAccent:        {},
	key.Comma:              {},
	key.FullStop:           {},
	key.Slash:              {},
	key.CapsLock:           {},

	key.F1:  {},
	key.F2:  {},
	key.F3:  {},
	key.F4:  {},
	key.F5:  {},
	key.F6:  {},
	key.F7:  {},
	key.F8:  {},
	key.F9:  {},
	key.F10: {},
	key.F11: {},
	key.F12: {},

	key.Pause:         {},
	key.Insert:        {},
	key.Home:          {},
	key.PageUp:        {},
	key.DeleteForward: {},
	key.End:           {},
	key.PageDown:      {},

	key.RightArrow: {},
	key.LeftArrow:  {},
	key.DownArrow:  {},
	key.UpArrow:    {},

	key.KeypadNumLock:     {},
	key.KeypadSlash:       {},
	key.KeypadAsterisk:    {},
	key.KeypadHyphenMinus: {},
	key.KeypadPlusSign:    {},
	key.KeypadEnter:       {},
	key.Keypad1:           {},
	key.Keypad2:           {},
	key.Keypad3:           {},
	key.Keypad4:           {},
	key.Keypad5:           {},
	key.Keypad6:           {},
	key.Keypad7:           {},
	key.Keypad8:           {},
	key.Keypad9:           {},
	key.Keypad0:           {},
	key.KeypadFullStop:    {},
	key.KeypadEqualSign:   {},

	key.F13: {},
	key.F14: {},
	key.F15: {},
	key.F16: {},
	key.F17: {},
	key.F18: {},
	key.F19: {},
	key.F20: {},
	key.F21: {},
	key.F22: {},
	key.F23: {},
	key.F24: {},

	key.Help: {},

	key.Mute:       {},
	key.VolumeUp:   {},
	key.VolumeDown: {},

	key.LeftControl:  {},
	key.LeftShift:    {},
	key.LeftAlt:      {},
	key.LeftGUI:      {},
	key.RightControl: {},
	key.RightShift:   {},
	key.RightAlt:     {},
	key.RightGUI:     {},
}
