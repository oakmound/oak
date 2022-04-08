package common

import (
	"image"
	"image/draw"

	"github.com/oakmound/oak/v3/shiny/screen"
)

type Image image.RGBA

func (i *Image) Size() image.Point {
	return i.Rect.Max
}

func (i *Image) Bounds() image.Rectangle {
	return i.Bounds()
}

func (Image) Release() {}

func (i *Image) RGBA() *image.RGBA {
	return (*image.RGBA)(i)
}

func (i *Image) Upload(dp image.Point, src screen.Image, sr image.Rectangle) {
	draw.Draw((*image.RGBA)(i), sr.Sub(sr.Min).Add(dp), src.RGBA(), sr.Min, draw.Src)
}

func NewImage(size image.Point) *Image {
	img := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))
	return (*Image)(img)
}
