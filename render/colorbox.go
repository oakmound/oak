package render

import (
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
	"image/draw"
)

type ColorBox struct {
	x, y  float64
	r     *image.RGBA
	layer int
}

func NewColorBox(w, h int, c color.Color) *ColorBox {
	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, image.NewUniform(c), image.Point{0, 0}, draw.Src)
	return &ColorBox{
		0.0,
		0.0,
		rgba,
		0,
	}
}

func (cb *ColorBox) GetRGBA() *image.RGBA {
	return cb.r
}

func (cb *ColorBox) Draw(buff screen.Buffer) {
	shinyDraw(buff, cb.r, int(cb.x), int(cb.y))
	draw.Draw(buff.RGBA(), buff.Bounds(),
		cb.r, image.Point{int(cb.x),
			int(cb.y)}, draw.Over)
}

func (cb *ColorBox) ShiftX(x float64) {
	cb.x += x
}
func (cb *ColorBox) ShiftY(y float64) {
	cb.y += y
}

func (cb *ColorBox) GetLayer() int {
	return cb.layer
}

func (cb *ColorBox) SetLayer(l int) {
	cb.layer = l
}

func (cb *ColorBox) UnDraw() {
	cb.layer = -1
}

func (cb *ColorBox) SetPos(x, y float64) {
	cb.x = x
	cb.y = y
}
