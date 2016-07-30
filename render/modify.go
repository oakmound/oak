package render

import (
	"github.com/disintegration/gift"
	"image"
	"image/color"
	"math"
)

// Modifications enum
const (
	F_FlipX = iota
	F_FlipY
	F_ApplyColor
	F_FillMask
	F_ApplyMask
	F_Rotate
	F_Scale
)

// Modifiables are Renderables that have functions to change their
// underlying image.
// This may be replaced with the gift library down the line
type Modifiable interface {
	Renderable
	FlipX()
	FlipY()
	ApplyColor(c color.Color)
	Copy() Modifiable
	FillMask(img image.RGBA)
	ApplyMask(img image.RGBA)
	Rotate(degrees int)
	Scale(xRatio float64, yRatio float64)
}

// FlipX returns a new rgba which is flipped
// over the horizontal axis.
func FlipX(rgba *image.RGBA) *image.RGBA {
	bounds := rgba.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y
	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			newRgba.Set(x, y, rgba.At(w-x, y))
		}
	}
	return newRgba
}

// FlipY returns a new rgba which is flipped
// over the vertical axis.
func FlipY(rgba *image.RGBA) *image.RGBA {
	bounds := rgba.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y
	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			newRgba.Set(x, y, rgba.At(x, h-y))
		}
	}
	return newRgba
}

// Apply color mixes a color into the rgba values of an image
// and returns that new rgba.
func ApplyColor(rgba *image.RGBA, c color.Color) *image.RGBA {
	r1, g1, b1, a1 := c.RGBA()
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

// FillMask replaces alpha 0 pixels in an RGBA with corresponding
// pixels in a second RGBA.
func FillMask(rgba *image.RGBA, img image.RGBA) *image.RGBA {
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

// ApplyMask mixes the rgba values of two images, according to
// their alpha levels, and returns that as a new rgba.
func ApplyMask(rgba *image.RGBA, img image.RGBA) *image.RGBA {
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

// Rotate returns a rotated rgba.
func Rotate(rgba *image.RGBA, degrees int) *image.RGBA {
	filter := gift.New(
		gift.Rotate(float32(degrees), color.Black, gift.CubicInterpolation))
	dst := image.NewRGBA(filter.Bounds(rgba.Bounds()))
	filter.Draw(dst, rgba)
	return dst

}

// Scale returns a scaled rgba.
func Scale(rgba *image.RGBA, xRatio float64, yRatio float64) *image.RGBA {
	bounds := rgba.Bounds()
	w := int(math.Floor(float64(bounds.Max.X) * xRatio))
	h := int(math.Floor(float64(bounds.Max.Y) * yRatio))
	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			newRgba.Set(x, y, rgba.At(int(math.Floor(float64(x)/xRatio)), int(math.Floor(float64(y)/yRatio))))
		}
	}
	return newRgba
}

func round(f float64) int {
	if f < -0.5 {
		return int(f - 0.5)
	}
	if f > 0.5 {
		return int(f + 0.5)
	}
	return 0
}
