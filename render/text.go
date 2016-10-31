package render

import (
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
)

type Text struct {
	Point
	Layered
	text string
	d_p  *font.Drawer
}

func NewStaticText(str string, x, y float64) *Text {
	return &Text{
		Point: Point{
			x,
			y,
		},
		text: str,
		d_p:  static_d,
	}
}

func NewText(str string, x, y float64) *Text {
	return &Text{
		Point: Point{
			x,
			y,
		},
		text: str,
		d_p:  d,
	}
}

func (t_p *Text) GetRGBA() *image.RGBA {
	return nil
}

func (t_p *Text) Draw(buff draw.Image) {
	t_p.d_p.Dot = fixed.P(int(t_p.X), int(t_p.Y))
	t_p.d_p.DrawString(t_p.text)
}

// Center will shift the text so that the existing leftmost point
// where the text sits becomes the center of the new text.
func (t_p *Text) Center() {
	textWidth := t_p.d_p.MeasureString(t_p.text).Round()
	t_p.ShiftX(float64(-textWidth / 2))
}

func (t_p *Text) SetText(str string) {
	t_p.text = str
	SetDirty(t_p.X, t_p.Y)
}

func (t *Text) String() string {
	return "Text[" + t.text + "]"
}

func StaticDrawText(str string, x, y int) {
	static_d.Dot = fixed.P(x, y)
	static_d.DrawString(str)
}

func DrawText(str string, x, y int) {
	d.Dot = fixed.P(x, y)
	d.DrawString(str)
}
