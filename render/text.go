package render

import (
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
)

type Text struct {
	LayeredPoint
	text string
	d    *Font
}

func (f *Font) NewText(str string, x, y float64) *Text {
	return &Text{
		LayeredPoint: LayeredPoint{
			Point: Point{
				X: x,
				Y: y,
			},
		},
		text: str,
		d:    f,
	}
}

func (t *Text) GetRGBA() *image.RGBA {
	return nil
}

func (t *Text) Draw(buff draw.Image) {
	t.d.Drawer.Dst = buff
	t.d.Drawer.Dot = fixed.P(int(t.X), int(t.Y))
	t.d.DrawString(t.text)
}

// Center will shift the text so that the existing leftmost point
// where the text sits becomes the center of the new text.
func (t *Text) Center() {
	textWidth := t.d.MeasureString(t.text).Round()
	t.ShiftX(float64(-textWidth / 2))
}

func (t *Text) SetText(str string) {
	t.text = str
}

func (t *Text) String() string {
	return "Text[" + t.text + "]"
}
