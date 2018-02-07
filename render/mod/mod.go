package mod

import (
	"image"
	"image/color"
	"math"

	"github.com/disintegration/gift"
)

// A Mod takes an image and returns that image transformed in some way.
type Mod func(image.Image) *image.RGBA

// A Transform is a longer name for writing Mod
type Transform = Mod

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
		newW := w - xOff + 1
		newH := h - yOff + 1
		if newW <= 0 || newH <= 0 {
			newW = 0
			newH = 0
		}
		out := image.NewRGBA(image.Rect(0, 0, newW, newH))
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
