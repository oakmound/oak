package render

import (
	"image/draw"
	"strconv"

	"bitbucket.org/oakmoundstudio/oak/physics"
	"golang.org/x/image/math/fixed"
)

type Text struct {
	LayeredPoint
	text Stringer
	d    *Font
}

type Stringer interface {
	String() string
}

func (f *Font) NewText(str Stringer, x, y float64) *Text {
	return &Text{
		LayeredPoint: LayeredPoint{
			Vector: physics.NewVector(x, y),
		},
		text: str,
		d:    f,
	}
}

func (t *Text) DrawOffset(buff draw.Image, xOff, yOff float64) {
	t.d.Drawer.Dst = buff
	t.d.Drawer.Dot = fixed.P(int(t.X()+xOff), int(t.Y()+yOff))
	t.d.DrawString(t.text.String())
}

func (t *Text) Draw(buff draw.Image) {
	t.d.Drawer.Dst = buff
	t.d.Drawer.Dot = fixed.P(int(t.X()), int(t.Y()))
	t.d.DrawString(t.text.String())
}

func (t *Text) SetFont(f *Font) {
	t.d = f
}

// Center will shift the text so that the existing leftmost point
// where the text sits becomes the center of the new text.
func (t *Text) Center() {
	textWidth := t.d.MeasureString(t.text.String()).Round()
	t.ShiftX(float64(-textWidth / 2))
}

func (t *Text) SetText(str Stringer) {
	t.text = str
}

func (t *Text) SetString(str string) {
	t.text = stringStringer(str)
}

func (t *Text) SetInt(i int) {
	t.text = stringStringer(strconv.Itoa(i))
}

func (t *Text) SetIntP(i *int) {
	t.text = stringerIntPointer{i}
}

// Todo: more SetX methods like float, floatP

func (t *Text) String() string {
	return "Text[" + t.text.String() + "]"
}

func (t *Text) Wrap(charLimit int, v float64) []*Text {
	st := t.text.String()
	out := make([]*Text, (len(st)/charLimit)+1)
	start := 0
	vertical := 0.0
	for i := range out {
		if start+charLimit <= len(st) {
			out[i] = t.d.NewStrText(st[start:start+charLimit], t.X(), t.Y()+vertical)
		} else {
			out[i] = t.d.NewStrText(st[start:], t.X(), t.Y()+vertical)
		}
		start += charLimit
		vertical += v
	}
	return out
}

type stringerIntPointer struct {
	v *int
}

func (sip stringerIntPointer) String() string {
	return strconv.Itoa(*sip.v)
}

// NewIntText wraps the given int pointer in a stringer interface
func (f *Font) NewIntText(str *int, x, y float64) *Text {
	return f.NewText(stringerIntPointer{str}, x, y)
}

type stringStringer string

func (ss stringStringer) String() string {
	return string(ss)
}

func (f *Font) NewStrText(str string, x, y float64) *Text {
	return f.NewText(stringStringer(str), x, y)
}
