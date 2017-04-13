package render

import (
	"image"
	"image/draw"

	"bitbucket.org/oakmoundstudio/oak/physics"
	"golang.org/x/image/math/fixed"
)

type IFText struct {
	LayeredPoint
	text Stringer
	d    *Font
}

type Stringer interface {
	String() string
}

func (f *Font) NewInterfaceText(str Stringer, x, y float64) *IFText {
	return &IFText{
		LayeredPoint: LayeredPoint{
			Vector: physics.Vector{
				X: x,
				Y: y,
			},
		},
		text: str,
		d:    f,
	}
}

func (t *IFText) GetRGBA() *image.RGBA {
	return nil
}

func (t *IFText) DrawOffset(buff draw.Image, xOff, yOff float64) {
	t.d.Drawer.Dst = buff
	t.d.Drawer.Dot = fixed.P(int(t.X+xOff), int(t.Y+yOff))
	t.d.DrawString(t.text.String())
}

func (t *IFText) Draw(buff draw.Image) {
	t.d.Drawer.Dst = buff
	t.d.Drawer.Dot = fixed.P(int(t.X), int(t.Y))
	t.d.DrawString(t.text.String())
}

// Center will shift the text so that the existing leftmost point
// where the text sits becomes the center of the new text.
func (t *IFText) Center() {
	textWidth := t.d.MeasureString(t.text.String()).Round()
	t.ShiftX(float64(-textWidth / 2))
}

func (t *IFText) SetText(str Stringer) {
	t.text = str
}

func (t *IFText) String() string {
	return "Text[" + t.text.String() + "]"
}
