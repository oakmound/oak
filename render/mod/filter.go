package mod

import (
	"image"
	"image/color"
)

// A Filter modifies an input image in place. This is useful notably for modifying
// a screen buffer, as they will refuse to be modified in any other way. This cannot
// change the dimensions of the underlying image.
type Filter func(*image.RGBA)

// ConformToPalleteFilter is not a modification, but acts like ConformToPallete
// without allocating a new *image.RGBA
func ConformToPalleteFilter(p color.Palette) Filter {
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

// InPlace converts a Mod to a Filter.
func InPlace(m Mod) Filter {
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
