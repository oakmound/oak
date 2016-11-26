package render

import (
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"strconv"
)

type IntText struct {
	LayeredPoint
	text *int
	d    *Font
}

func (f *Font) NewIntText(str *int, x, y float64) *IntText {
	return &IntText{
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

func (t *IntText) GetRGBA() *image.RGBA {
	return nil
}

func (t *IntText) DrawOffset(buff draw.Image, xOff, yOff float64) {
	t.d.Drawer.Dst = buff
	t.d.Drawer.Dot = fixed.P(int(t.X+xOff), int(t.Y+yOff))
	t.d.DrawString(strconv.Itoa(*t.text))
}

func (t *IntText) Draw(buff draw.Image) {
	// We need to benchmark if this buff replacement is slow
	t.d.Drawer.Dst = buff
	t.d.Drawer.Dot = fixed.P(int(t.X), int(t.Y))
	t.d.DrawString(strconv.Itoa(*t.text))
}

// Center will shift the text so that the existing leftmost point
// where the text sits becomes the center of the new text.
func (t *IntText) Center() {
	textWidth := t.d.MeasureString(strconv.Itoa(*t.text)).Round()
	t.ShiftX(float64(-textWidth / 2))
}

func (t *IntText) SetText(str *int) {
	t.text = str
}
