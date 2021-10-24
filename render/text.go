package render

import (
	"fmt"
	"image/draw"
	"strconv"

	"github.com/oakmound/oak/v3/alg"
	"golang.org/x/image/math/fixed"
)

// A Text is a renderable that represents some text to print on screen
type Text struct {
	LayeredPoint
	text fmt.Stringer
	d    *Font
}

// NewStringerText creates a renderable text component that will draw the string
// provided by the given stringer each frame.
func (f *Font) NewStringerText(str fmt.Stringer, x, y float64) *Text {
	return &Text{
		LayeredPoint: NewLayeredPoint(x, y, 0),
		text:         str,
		d:            f.Copy(),
	}
}

type stringerIntPointer struct {
	v *int
}

func (sip stringerIntPointer) String() string {
	return strconv.Itoa(*sip.v)
}

// NewIntText wraps the given int pointer in a stringer interface
func (f *Font) NewIntText(str *int, x, y float64) *Text {
	return f.NewStringerText(stringerIntPointer{str}, x, y)
}

type stringStringer string

func (ss stringStringer) String() string {
	return string(ss)
}

// NewText creates a renderable text component with the given string body
func (f *Font) NewText(str string, x, y float64) *Text {
	return f.NewStringerText(stringStringer(str), x, y)
}

type stringPtrStringer struct {
	s *string
}

func (sp stringPtrStringer) String() string {
	if sp.s == nil {
		return "nil"
	}
	return string(*sp.s)
}

// NewStrPtrText creates a renderable text component with a body matching
// and updating to match the content behind the provided string pointer
func (f *Font) NewStrPtrText(str *string, x, y float64) *Text {
	return f.NewStringerText(stringPtrStringer{str}, x, y)
}

func (t *Text) drawWithFont(buff draw.Image, xOff, yOff float64, fnt *Font) {
	fnt.Drawer.Dst = buff
	sz := 12.0
	if t.d.gen.Size != 0 {
		sz = t.d.gen.Size
	}
	fnt.Drawer.Dot = fixed.P(int(t.X()+xOff), int(t.Y()+yOff)+int(sz))
	fnt.drawString(t.text.String())
}

// Draw for a text draws the text at its layeredPoint position
func (t *Text) Draw(buff draw.Image, xOff, yOff float64) {
	t.drawWithFont(buff, xOff, yOff, t.d)
}

// SetFont sets the drawer which renders the text each frame
func (t *Text) SetFont(f *Font) {
	t.d = f
}

// GetDims reports the width and height of a text renderable
func (t *Text) GetDims() (int, int) {
	// BUG: reported height is too low, test this impl:
	// bounds, adv := t.d.BoundString(t.text.String())
	// return adv.Round(), bounds.Max.Y.Round()
	textWidth := t.d.MeasureString(t.text.String()).Round()
	return textWidth, alg.RoundF64(t.d.gen.Size)
}

// Center will shift the text so that the existing leftmost point
// where the text sits becomes the center of the new text.
func (t *Text) Center() {
	textWidth := t.d.MeasureString(t.text.String()).Round()
	t.ShiftX(float64(-textWidth / 2))
}

// SetString accepts a string itself as the stringer to be written
func (t *Text) SetString(str string) {
	t.text = stringStringer(str)
}

// SetStringPtr accepts a string pointer as the stringer to be written
func (t *Text) SetStringPtr(str *string) {
	t.text = stringPtrStringer{str}
}

// SetStringer accepts an fmt.Stringer to write
func (t *Text) SetStringer(s fmt.Stringer) {
	t.text = s
}

//SetInt takes and converts the input integer to a string to write
func (t *Text) SetInt(i int) {
	t.text = stringStringer(strconv.Itoa(i))
}

// SetIntPtr takes in an integer pointer that will draw the integer
// behind the pointer, in base 10, each frame
func (t *Text) SetIntPtr(i *int) {
	t.text = stringerIntPointer{i}
}

// StringLiteral returns what text is currently rendering.
func (t *Text) StringLiteral() string {
	return t.text.String()
}

// Wrap returns the input text split into a list of texts
// spread vertically, splitting after each charLimit is reached.
// the input vertInc is how much each text in the slice will differ by
// in y value
func (t *Text) Wrap(charLimit int, vertInc float64) []*Text {
	st := t.text.String()
	outlen := len(st) / charLimit
	if len(st)%charLimit != 0 {
		outlen++
	}
	out := make([]*Text, outlen)
	start := 0
	vertical := 0.0
	for i := range out {
		if start+charLimit <= len(st) {
			out[i] = t.d.NewText(st[start:start+charLimit], t.X(), t.Y()+vertical)
		} else {
			out[i] = t.d.NewText(st[start:], t.X(), t.Y()+vertical)
		}
		start += charLimit
		vertical += vertInc
	}
	return out
}

// ToSprite converts this text into a sprite, so that it is no longer
// modifiable in terms of its text content, but is modifiable in terms
// of mod.Transform or mod.Filter.
func (t *Text) ToSprite() *Sprite {
	tmpFnt := t.d.Copy()
	width := tmpFnt.MeasureString(t.text.String()).Round()
	height := tmpFnt.bounds.Max.Y()
	s := NewEmptySprite(t.X(), t.Y(), width, height+5)
	t.drawWithFont(s.GetRGBA(), -t.X(), -t.Y(), tmpFnt)
	return s
}
