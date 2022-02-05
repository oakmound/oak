//go:build js
// +build js

package jsdriver

import "image"

type imageImpl struct {
	screen *screenImpl
	size   image.Point
	rgba   *image.RGBA
}

func (ii imageImpl) Size() image.Point {
	return ii.size
}

func (ii imageImpl) Bounds() image.Rectangle {
	return image.Rect(0, 0, ii.size.X, ii.size.Y)
}

func (imageImpl) Release() {}

func (ii imageImpl) RGBA() *image.RGBA {
	return ii.rgba
}
