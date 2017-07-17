package render

import (
	"fmt"
	"image/draw"
	"strconv"

	"github.com/oakmound/oak/physics"
	"golang.org/x/image/math/fixed"
)

// A Text is a renderable that represents some text to print on screen
type Text struct {
	LayeredPoint
	text fmt.Stringer
	d    *Font
}

// NewText takes in anything that has a String() function and returns a text
// object with the associated font and screen position
func (f *Font) NewText(str fmt.Stringer, x, y float64) *Text {
	return &Text{
		LayeredPoint: LayeredPoint{
			Vector: physics.NewVector(x, y),
		},
		text: str,
		d:    f,
	}
}

// DrawOffset for a text object draws the text at t.(X,Y) + (xOff,yOff)
func (t *Text) DrawOffset(buff draw.Image, xOff, yOff float64) {
	t.d.Drawer.Dst = buff
	t.d.Drawer.Dot = fixed.P(int(t.X()+xOff), int(t.Y()+yOff))
	t.d.DrawString(t.text.String())
}

// Draw for a text draws the text at its layeredPoint position
func (t *Text) Draw(buff draw.Image) {
	t.d.Drawer.Dst = buff
	t.d.Drawer.Dot = fixed.P(int(t.X()), int(t.Y()))
	t.d.DrawString(t.text.String())
}

// SetFont sets the drawer which renders the text each frame
func (t *Text) SetFont(f *Font) {
	t.d = f
}

// Center will shift the text so that the existing leftmost point
// where the text sits becomes the center of the new text.
func (t *Text) Center() {
	textWidth := t.d.MeasureString(t.text.String()).Round()
	t.ShiftX(float64(-textWidth / 2))
}

// SetText sets the string to be written to a new stringer
func (t *Text) SetText(str fmt.Stringer) {
	t.text = str
}

// SetString accepts a string itself as the stringer to be written
func (t *Text) SetString(str string) {
	t.text = stringStringer(str)
}

//SetInt takes and converts the input integer to a string to write
func (t *Text) SetInt(i int) {
	t.text = stringStringer(strconv.Itoa(i))
}

// SetIntP takes in an integer pointer that will be drawn at whatever
// the value is behind the pointer when it is drawn
func (t *Text) SetIntP(i *int) {
	t.text = stringerIntPointer{i}
}

// Todo: more SetX methods like float, floatP

func (t *Text) String() string {
	return "Text[" + t.text.String() + "]"
}

// Wrap returns the input text split into a list of texts
// spread vertically, splitting after each charLimit is reached.
// the input vertInc is how much each text in the slice will differ by
// in y value
func (t *Text) Wrap(charLimit int, vertInc float64) []*Text {
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
		vertical += vertInc
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

// NewStrText is a helper to take in a string instead of a stringer for NewText
func (f *Font) NewStrText(str string, x, y float64) *Text {
	return f.NewText(stringStringer(str), x, y)
}
