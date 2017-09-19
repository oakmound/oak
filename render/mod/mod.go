package mod

import (
	"image"
	"image/color"
	"math"

	"github.com/disintegration/gift"
)

// A Mod takes an image and returns that image transformed in some way.
type Mod func(image.Image) *image.RGBA

// And chains together multiple Mods into a single Mod
func And(ms ...Mod) Mod {
	return func(rgba image.Image) *image.RGBA {
		rgba2 := ms[0](rgba)
		for i := 1; i < len(ms); i++ {
			rgba2 = ms[i](rgba2)
		}
		return rgba2
	}
}

// Scale returns a scaled rgba.
func Scale(xRatio, yRatio float64) Mod {
	return func(rgba image.Image) *image.RGBA {
		bounds := rgba.Bounds()
		w := int(math.Floor(float64(bounds.Max.X) * xRatio))
		h := int(math.Floor(float64(bounds.Max.Y) * yRatio))
		filter := gift.New(
			gift.Resize(w, h, gift.CubicResampling))
		dst := image.NewRGBA(filter.Bounds(rgba.Bounds()))
		filter.Draw(dst, rgba)
		return dst
	}
}

// Fade reduces the alpha of an image
func Fade(alpha int) Mod {
	return func(rgba image.Image) *image.RGBA {
		bounds := rgba.Bounds()
		a2 := uint32(alpha * 257)
		w := bounds.Max.X
		h := bounds.Max.Y
		newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				r, g, b, a := rgba.At(x, y).RGBA()
				var a3 uint16
				if a2 > a {
					a3 = 0
				} else {
					a3 = uint16(a - a2)
				}
				tmp := color.NRGBA64{
					uint16(r),
					uint16(g),
					uint16(b),
					a3}
				newRgba.Set(x, y, tmp)
			}
		}
		return newRgba
	}
}

// ApplyColor mixes a color into the rgba values of an image
// and returns that new rgba.
func ApplyColor(c color.Color) Mod {
	return func(rgba image.Image) *image.RGBA {

		// u := image.NewUniform(c)
		// bounds := rgba.Bounds()
		// img := image.NewRGBA(bounds)
		// draw.Draw(img, bounds, u, bounds.Min, draw.Src)

		// return blend.Normal(rgba, img)
		r1, g1, b1, a1 := c.RGBA()
		// filter := gift.New(
		// 	gift.ColorBalance(float32(r1*(255/100)), float32(g1*(255/100)), float32(b1*(255/100))))
		// dst := image.NewRGBA(filter.Bounds(rgba.Bounds()))
		// filter.Draw(dst, rgba)
		// return dst

		bounds := rgba.Bounds()
		w := bounds.Max.X
		h := bounds.Max.Y
		newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				r2, g2, b2, a2 := rgba.At(x, y).RGBA()
				a3 := a1 + a2
				if a2 == 0 {
					newRgba.Set(x, y, color.RGBA{0, 0, 0, 0})
					continue
				}
				tmp := color.RGBA64{
					uint16(((a1 * r1) + (a2 * r2)) / a3),
					uint16(((a1 * g1) + (a2 * g2)) / a3),
					uint16(((a1 * b1) + (a2 * b2)) / a3),
					uint16(a2)}
				newRgba.Set(x, y, tmp)
			}
		}
		return newRgba
	}
}

// FillMask replaces alpha 0 pixels in an RGBA with corresponding
// pixels in a second RGBA.
func FillMask(img image.RGBA) Mod {
	return func(rgba image.Image) *image.RGBA {
		// Instead of static color it just two buffers melding
		bounds := rgba.Bounds()
		w := bounds.Max.X
		h := bounds.Max.Y
		newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				newRgba.Set(x, y, rgba.At(x, y))
			}
		}
		bounds = img.Bounds()
		w = bounds.Max.X
		h = bounds.Max.Y
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

				newRgba.Set(x, y, tmp)
			}
		}
		return newRgba
	}
}

// ApplyMask mixes the rgba values of two images, according to
// their alpha levels, and returns that as a new rgba.
func ApplyMask(img image.RGBA) Mod {
	return func(rgba image.Image) *image.RGBA {
		// Instead of static color it just two buffers melding
		bounds := rgba.Bounds()
		w := bounds.Max.X
		h := bounds.Max.Y
		newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				newRgba.Set(x, y, rgba.At(x, y))
			}
		}
		bounds = img.Bounds()
		w = bounds.Max.X
		h = bounds.Max.Y
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				r1, g1, b1, a1 := img.At(x, y).RGBA()
				r2, g2, b2, a2 := rgba.At(x, y).RGBA()

				var tmp color.RGBA64

				a3 := a1 + a2
				if a3 == 0 {
					tmp = color.RGBA64{
						0, 0, 0, 0,
					}
					newRgba.Set(x, y, tmp)
					continue
				}

				tmp = color.RGBA64{
					uint16(((a1 * r1) + (a2 * r2)) / a3),
					uint16(((a1 * g1) + (a2 * g2)) / a3),
					uint16(((a1 * b1) + (a2 * b2)) / a3),
					uint16(math.Max(float64(a1), float64(a2)))}

				newRgba.Set(x, y, tmp)
			}
		}
		return newRgba
	}
}

// TrimColor will trim inputs so that any rows or columns where each pixel is
// less than or equal to the input color are removed. This will change the dimensions
// of the image.
func TrimColor(trimUnder color.Color) Mod {
	r, g, b, a := trimUnder.RGBA()
	return func(rgba image.Image) *image.RGBA {
		bounds := rgba.Bounds()
		w := bounds.Max.X
		h := bounds.Max.Y
		xOff := 0
		yOff := 0
	trimouter1:
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				c := rgba.At(x, y)
				r2, g2, b2, a2 := c.RGBA()
				if colorLess(r, r2, g, g2, b, b2, a, a2) {
					continue
				}
				break trimouter1
			}
			xOff++
		}
	trimouter2:
		for x := w; x >= 0; x-- {
			for y := 0; y < h; y++ {
				c := rgba.At(x, y)
				r2, g2, b2, a2 := c.RGBA()
				if colorLess(r, r2, g, g2, b, b2, a, a2) {
					continue
				}
				break trimouter2
			}
			w--
		}
	trimouter3:
		for y := h; y >= 0; y-- {
			for x := 0; x < w; x++ {
				c := rgba.At(x, y)
				r2, g2, b2, a2 := c.RGBA()
				if colorLess(r, r2, g, g2, b, b2, a, a2) {
					continue
				}
				break trimouter3
			}
			h--
		}
	trimouter4:
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				c := rgba.At(x, y)
				r2, g2, b2, a2 := c.RGBA()
				if colorLess(r, r2, g, g2, b, b2, a, a2) {
					continue
				}
				break trimouter4
			}
			yOff++
		}
		out := image.NewRGBA(image.Rect(0, 0, w-xOff+1, h-yOff+1))
		for x := xOff; x <= w; x++ {
			for y := yOff; y <= h; y++ {
				c := rgba.At(x, y)
				out.Set(x-xOff, y-yOff, c)
			}
		}
		return out
	}
}

func colorLess(r, r2, g, g2, b, b2, a, a2 uint32) bool {
	return r2 <= r && g2 <= g && b2 <= b && a2 <= a
}

// ConformToPallete modifies the input image so that it's colors all fall
// in the input color palette.
func ConformToPallete(p color.Model) Mod {
	return func(rgba image.Image) *image.RGBA {
		bounds := rgba.Bounds()
		w := bounds.Max.X
		h := bounds.Max.Y
		newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				newRgba.Set(x, y, p.Convert(rgba.At(x, y)))
			}
		}
		return newRgba
	}
}

// Zoom zooms into a position on the input image.
// The position is determined by the input percentages, and how far the zoom
// is deep depends on the input zoom level-- 2.0 would quarter the number of
// unique pixels from the input to the output.
func Zoom(xPerc, yPerc, zoom float64) func(rgba image.Image) *image.RGBA {
	return func(rgba image.Image) *image.RGBA {
		bounds := rgba.Bounds()
		w := float64(bounds.Max.X)
		h := float64(bounds.Max.Y)
		newRgba := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
		newW := w / zoom
		newH := h / zoom
		minX := (w - newW) * xPerc
		minY := (h - newH) * yPerc
		for x := 0.0; x < w; x++ {
			for y := 0.0; y < h; y++ {
				x2 := int(((x * xPerc) / (zoom * xPerc)) + minX)
				y2 := int(((y * yPerc) / (zoom * yPerc)) + minY)
				newRgba.Set(int(x), int(y), rgba.At(x2, y2))
			}
		}
		return newRgba
	}
}
