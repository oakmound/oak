package common

import "image"

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

func NewImage(size image.Point) *Image {
	img := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))
	return (*Image)(img)
}
