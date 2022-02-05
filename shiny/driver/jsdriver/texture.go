//go:build js
// +build js

package jsdriver

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/oakmound/oak/v3/shiny/screen"
)

type textureImpl struct {
	screen *screenImpl
	size   image.Point
	rgba   *image.RGBA
}

func (ti *textureImpl) Size() image.Point {
	return ti.size
}

func (ti *textureImpl) Bounds() image.Rectangle {
	return image.Rect(0, 0, ti.size.X, ti.size.Y)
}

func (ti *textureImpl) Upload(dp image.Point, src screen.Image, sr image.Rectangle) {
	rgba := src.RGBA()
	ti.rgba = rgba
}
func (*textureImpl) Fill(dr image.Rectangle, src color.Color, op draw.Op) {}
func (*textureImpl) Release()                                             {}
