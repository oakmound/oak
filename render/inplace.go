package render

import (
	"image"
	"image/color"
)

// An InPlaceMod operates like a Modification, but it will modify the original rgba
// contents as opposed to returning a new transformed rgba.
type InPlaceMod func(*image.RGBA)

// ConformToPalleteInPlace is not a modification, but acts like ConformToPallete
// without allocating a new *image.RGBA
func ConformToPalleteInPlace(p color.Palette) func(rgba *image.RGBA) {
	return func(rgba *image.RGBA) {
		bounds := rgba.Bounds()
		w := bounds.Max.X
		h := bounds.Max.Y
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				rgba.Set(x, y, p.Convert(rgba.At(x, y)))
			}
		}
	}
}

// InPlace converts a Modification to an InPlaceMod.
func InPlace(m Modification) InPlaceMod {
	return func(rgba *image.RGBA) {
		rgba2 := m(rgba)
		bounds := rgba.Bounds()
		w := bounds.Max.X
		h := bounds.Max.Y
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				rgba.Set(x, y, rgba2.At(x, y))
			}
		}
	}
}
