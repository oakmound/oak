package render

import (
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"strconv"
)

type IntText struct {
	LayeredPoint
	text *int
	d_p  *font.Drawer
}

func NewIntText(str *int, x, y float64) *IntText {
	return &IntText{
		LayeredPoint: LayeredPoint{
			Point: Point{
				X: x,
				Y: y,
			},
		},
		text: str,
		d_p:  d,
	}
}

func NewStaticIntText(str *int, x, y float64) *IntText {
	return &IntText{
		LayeredPoint: LayeredPoint{
			Point: Point{
				X: x,
				Y: y,
			},
		},
		text: str,
		d_p:  static_d,
	}
}

func (t_p *IntText) GetRGBA() *image.RGBA {
	return nil
}

func (t_p *IntText) Draw(buff draw.Image) {
	t_p.d_p.Dot = fixed.P(int(t_p.X), int(t_p.Y))
	t_p.d_p.DrawString(strconv.Itoa(*t_p.text))
}

// Center will shift the text so that the existing leftmost point
// where the text sits becomes the center of the new text.
func (t_p *IntText) Center() {
	textWidth := t_p.d_p.MeasureString(strconv.Itoa(*t_p.text)).Round()
	t_p.ShiftX(float64(-textWidth / 2))
}

func (t_p *IntText) SetText(str *int) {
	t_p.text = str
}
