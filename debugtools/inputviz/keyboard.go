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

type LayoutKey interface {
	Pos() LayoutPosition
}

type LayoutPosition struct {
	Key    string
	Gap    bool
	Width  float64
	Height float64
	Row    float64
	Col    float64
}

type gap float64

func (g gap) Pos() LayoutPosition {
	return LayoutPosition{
		Gap:   true,
		Width: float64(g),
	}
}

type standardKey string

func (s standardKey) Pos() LayoutPosition {
	return LayoutPosition{
		Key:    string(s),
		Width:  1,
		Height: 1,
	}
}

type wideKey struct {
	k string
	w float64
}

func (w wideKey) Pos() LayoutPosition {
	return LayoutPosition{
		Key:    string(w.k),
		Width:  w.w,
		Height: 1,
	}
}

type tallKey struct {
	k string
	h float64
}

func (h tallKey) Pos() LayoutPosition {
	return LayoutPosition{
		Key:    string(h.k),
		Width:  1,
		Height: h.h,
	}
}

type LayoutQWERTY struct {
	Bounds    floatgeom.Rect2
	layoutMap map[string]LayoutPosition
}

func (l *LayoutQWERTY) init() {
	if l.layoutMap != nil {
		return
	}
	type sk = standardKey

	l.layoutMap = make(map[string]LayoutPosition)
	qwertyRows := [][]LayoutKey{
		{sk(key.EscapeStr), gap(1), sk(key.F1Str), sk(key.F2Str), sk(key.F3Str), sk(key.F4Str), gap(.5), sk(key.F5Str), sk(key.F6Str), sk(key.F7Str), sk(key.F8Str), gap(.5), sk(key.F9Str), sk(key.F10Str), sk(key.F11Str), sk(key.F12Str), gap(2.1), sk(key.Pause)},
		{sk(key.GraveAccentStr), sk(key.OneStr), sk(key.TwoStr), sk(key.ThreeStr), sk(key.FourStr), sk(key.FiveStr), sk(key.SixStr), sk(key.SevenStr), sk(key.EightStr), sk(key.NineStr), sk(key.ZeroStr), sk(key.HyphenMinusStr), sk(key.EqualSignStr), wideKey{key.DeleteBackspaceStr, 2.0}, gap(.1), sk(key.InsertStr), sk(key.HomeStr), sk(key.PageUpStr), gap(.1), sk(key.KeypadNumLockStr), sk(key.KeypadSlashStr), sk(key.KeypadAsteriskStr), sk(key.KeypadHyphenMinus)},
		{wideKey{key.TabStr, 1.5}, sk(key.QStr), sk(key.WStr), sk(key.EStr), sk(key.RStr), sk(key.TStr), sk(key.YStr), sk(key.UStr), sk(key.IStr), sk(key.OStr), sk(key.PStr), sk(key.LeftSquareBracketStr), sk(key.RightSquareBracketStr), wideKey{key.BackslashStr, 1.5}, gap(.1), sk(key.DeleteForwardStr), sk(key.EndStr), sk(key.PageDownStr), gap(.1), sk(key.Keypad7Str), sk(key.Keypad8Str), sk(key.Keypad9Str), tallKey{key.KeypadPlusSignStr, 2}},
		{wideKey{key.CapsLockStr, 1.5}, sk(key.AStr), sk(key.SStr), sk(key.DStr), sk(key.FStr), sk(key.GStr), sk(key.HStr), sk(key.JStr), sk(key.KStr), sk(key.LStr), sk(key.SemicolonStr), sk(key.ApostropheStr), wideKey{key.ReturnEnterStr, 2.5}, gap(3.2), sk(key.Keypad4Str), sk(key.Keypad5Str), sk(key.Keypad6)},
		{wideKey{key.LeftShiftStr, 2.0}, sk(key.ZStr), sk(key.XStr), sk(key.CStr), sk(key.VStr), sk(key.BStr), sk(key.NStr), sk(key.MStr), sk(key.CommaStr), sk(key.FullStopStr), sk(key.SlashStr), wideKey{key.RightShiftStr, 3.0}, gap(1.1), sk(key.UpArrowStr), gap(1.1), sk(key.Keypad1Str), sk(key.Keypad2Str), sk(key.Keypad3Str), tallKey{key.KeypadEnterStr, 2.0}},
		{wideKey{key.LeftControlStr, 1.5}, sk(key.LeftGUIStr), wideKey{key.LeftAltStr, 1.5}, wideKey{key.SpacebarStr, 7.0}, wideKey{key.RightAltStr, 1.5}, sk(key.RightGUIStr), wideKey{key.RightControlStr, 1.5}, gap(.1), sk(key.LeftArrowStr), sk(key.DownArrowStr), sk(key.RightArrowStr), gap(.1), wideKey{key.Keypad0Str, 2.0}, sk(key.KeypadPeriod)},
	}
	rowFloats := []float64{0.0, 1.1, 2.1, 3.1, 4.1, 5.1}
	for row, cols := range qwertyRows {
		rf := rowFloats[row]
		cf := 0.0
		for _, v := range cols {
			ps := v.Pos()
			if ps.Key != "" {
				l.layoutMap[ps.Key] = LayoutPosition{
					Row:    rf,
					Col:    cf,
					Width:  ps.Width - .1,
					Height: ps.Height - .1,
				}
			}
			cf += ps.Width
		}
	}
}

func (l *LayoutQWERTY) KeyRect(k string) floatgeom.Rect2 {
	l.init()

	pos, ok := l.layoutMap[k]
	if !ok {
		return floatgeom.Rect2{}
	}
	row := pos.Row
	col := pos.Col
	width := pos.Width
	height := pos.Height

	w, h := l.Bounds.W(), l.Bounds.H()
	// max row = 5.1, drawn down to 6
	// max col = 21.2, drawn right to 22.1
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

	RenderCharacters bool
	Font             *render.Font

	event.CallerID
	ctx *scene.Context

	rs map[string]*render.Switch

	bindings []event.Binding
}

func (k *Keyboard) CID() event.CallerID {
	return k.CallerID
}

func (k *Keyboard) RenderAndListen(ctx *scene.Context, layer int) error {
	k.ctx = ctx
	k.CallerID = k.ctx.CallerMap.Register(k)

	if k.Rect.W() == 0 || k.Rect.H() == 0 {
		k.Rect.Max = k.Rect.Min.Add(floatgeom.Point2{320, 180})
	}
	if k.KeyboardLayout == nil {
		k.KeyboardLayout = &LayoutQWERTY{
			Bounds: k.Rect,
		}
	}
	if k.Colors == nil {
		k.Colors = defaultColors
	}
	if k.Font == nil {
		k.Font = render.DefaultFont()
	}

	k.rs = make(map[string]*render.Switch)

	for kv := range key.AllKeys {
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
		if k.RenderCharacters {
			x, y := rect.Min.X(), rect.Min.Y()
			txt := k.Font.NewText(kv, x, y)
			tw, th := txt.GetDims()
			xBuffer := rect.W() - float64(tw)
			yBuffer := rect.H() - float64(th)
			// Only render strings that will stay inside their boundaries
			if xBuffer >= 0 {
				txt.ShiftX(xBuffer / 2)
				txt.ShiftY(yBuffer / 2)
				if k.BaseLayer == -1 {
					ctx.DrawStack.Draw(txt, layer+1)
				} else {
					ctx.DrawStack.Draw(txt, k.BaseLayer, layer+1)
				}
			}
		}
		if k.BaseLayer == -1 {
			ctx.DrawStack.Draw(r, layer)
		} else {
			ctx.DrawStack.Draw(r, k.BaseLayer, layer)
		}
	}

	b1 := event.Bind(ctx, key.AnyDown, k, func(kb *Keyboard, ev key.Event) event.Response {
		btn := ev.Code.String()[4:]
		if kb.rs[btn] == nil {
			return 0
		}
		kb.rs[btn].Set("pressed")
		return 0
	})
	b2 := event.Bind(ctx, key.AnyUp, k, func(kb *Keyboard, ev key.Event) event.Response {
		btn := ev.Code.String()[4:]
		if kb.rs[btn] == nil {
			return 0
		}
		kb.rs[btn].Set("released")
		return 0
	})
	k.bindings = []event.Binding{b1, b2}
	return nil
}

func (k *Keyboard) Destroy() {
	for _, b := range k.bindings {
		b.Unbind()
	}
	for _, r := range k.rs {
		r.Undraw()
	}
}
