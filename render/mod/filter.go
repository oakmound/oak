package mod

import (
	"image"
	"image/color"
	"math"

	"github.com/oakmound/oak/v2/dlog"
)

// A Filter modifies an input image in place. This is useful notably for modifying
// a screen buffer, as they will refuse to be modified in any other way. This cannot
// change the dimensions of the underlying image.
type Filter func(*image.RGBA)

// AndFilter combines multiple filters into one.
func AndFilter(fs ...Filter) Filter {
	return func(rgba *image.RGBA) {
		for _, f := range fs {
			f(rgba)
		}
	}
}

// ConformToPallete is not a modification, but acts like ConformToPallete
// without allocating a new *image.RGBA
func ConformToPallete(p color.Model) Filter {
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

// Fade reduces the alpha of an image. It takes an alpha from 0-255.
func Fade(alpha int) Filter {
	return func(rgba *image.RGBA) {
		bounds := rgba.Bounds()
		a2 := uint32(alpha * 257)
		w := bounds.Max.X
		h := bounds.Max.Y
		var a3 uint16
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				r, g, b, a := rgba.At(x, y).RGBA()
				if a2 > a {
					a3 = 0
				} else {
					a3 = uint16(a - a2)
				}
				rgba.Set(x, y, color.NRGBA64{
					uint16(r),
					uint16(g),
					uint16(b),
					a3})
			}
		}
	}
}

// ApplyMask mixes the rgba values of two images, according to
// their alpha levels, and returns that as a new rgba.
func ApplyMask(img image.RGBA) Filter {
	return func(rgba *image.RGBA) {
		bounds := img.Bounds()
		w := bounds.Max.X
		h := bounds.Max.Y
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				r1, g1, b1, a1 := img.At(x, y).RGBA()
				r2, g2, b2, a2 := rgba.At(x, y).RGBA()

				a3 := a1 + a2
				if a3 == 0 {
					rgba.Set(x, y, color.RGBA64{0, 0, 0, 0})
					continue
				}

				rgba.Set(x, y, color.RGBA64{
					uint16(((a1 * r1) + (a2 * r2)) / a3),
					uint16(((a1 * g1) + (a2 * g2)) / a3),
					uint16(((a1 * b1) + (a2 * b2)) / a3),
					uint16(math.Max(float64(a1), float64(a2)))})
			}
		}
	}
}

// ApplyColor mixes a color into the rgba values of an image
// and returns that new rgba.
func ApplyColor(c color.Color) Filter {
	return func(rgba *image.RGBA) {

		r1, g1, b1, a1 := c.RGBA()
		bounds := rgba.Bounds()
		w := bounds.Max.X
		h := bounds.Max.Y
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				r2, g2, b2, a2 := rgba.At(x, y).RGBA()
				a3 := a1 + a2
				if a2 == 0 {
					rgba.Set(x, y, color.RGBA{0, 0, 0, 0})
					continue
				}
				rgba.Set(x, y, color.RGBA64{
					uint16(((a1 * r1) + (a2 * r2)) / a3),
					uint16(((a1 * g1) + (a2 * g2)) / a3),
					uint16(((a1 * b1) + (a2 * b2)) / a3),
					uint16(a2)})
			}
		}
	}
}

// FillMask replaces alpha 0 pixels in an RGBA with corresponding
// pixels in a second RGBA.
func FillMask(img image.RGBA) Filter {
	return func(rgba *image.RGBA) {
		bounds := img.Bounds()
		w := bounds.Max.X
		h := bounds.Max.Y
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				r1, g1, b1, a1 := rgba.At(x, y).RGBA()
				r2, g2, b2, a2 := img.At(x, y).RGBA()

				var tmp color.RGBA64

				if a1 == 0 {
					tmp = color.RGBA64{
						uint16(r2),
						uint16(g2),
						uint16(b2),
						uint16(a2),
					}
				} else {
					tmp = color.RGBA64{
						uint16(r1),
						uint16(g1),
						uint16(b1),
						uint16(a1),
					}
				}

				rgba.Set(x, y, tmp)
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

// StripOuterAlpha from the image given a source image and a alpha level to denote stripping.
// Note that this was implement for ease of implementation but not speed.
// We could use image lib or a real depth first search to do fewer checks but this is easier...
func StripOuterAlpha(m *image.RGBA, level int) Filter {
	l := uint8(level)
	return func(rgba *image.RGBA) {
		if m == nil {
			dlog.Warn("Invalid rgba provided to stripouteralpha")
			return
		}

		mr := m

		//get an image
		bounds := mr.Bounds()
		w := bounds.Max.X
		h := bounds.Max.Y

		// check downwards for the given level.x
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				r := mr.RGBAAt(x, y)
				if r.A <= l {
					// Treating as transparent
					rgba.Set(int(x), int(y), mr.At(x, y))
				} else {
					break
				}
			}
		}
		// check left to right for the given level.
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r := mr.RGBAAt(x, y)
				if r.A <= l {
					// Treating as transparent
					rgba.Set(int(x), int(y), mr.At(x, y))
				} else {
					break
				}
			}
		}

		// check bottom up for the given level.
		for x := 0; x < w; x++ {
			for y := h - 1; y >= 0; y-- {
				r := mr.RGBAAt(x, y)
				if r.A <= l {
					// Treating as transparent
					rgba.Set(int(x), int(y), mr.At(x, y))
				} else {
					break
				}
			}
		}
		// check right to left for the given level.
		for y := 0; y < h; y++ {
			for x := w - 1; x >= 0; x-- {
				r := mr.RGBAAt(x, y)
				if r.A <= l {
					// Treating as transparent
					rgba.Set(int(x), int(y), mr.At(x, y))
				} else {
					break
				}
			}
		}

		// apply image onto m

	}

}

// There is no function to convert a Filter to a Mod, to promote not doing so.
// Mods are significantly less efficient than Filters.
