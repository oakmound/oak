package render

import (
	// This file is being slowly converted to use gift over manual math and loops,
	// because our math / loops will be more likely to have (and have already had)
	// missable bugs.
	//"github.com/anthonynsimon/bild/blend"

	"github.com/disintegration/gift"

	"image"
	"image/color"
	//"image/draw"
	"math"
)

var (
	transparent = color.RGBA{0, 0, 0, 0}
)

// A Modification takes in an image buffer and returns a new image buffer
type Modification func(image.Image) *image.RGBA

// A Modifiable is a Renderable that has functions to change its
// underlying image.
// This may be replaced with the gift library down the line
type Modifiable interface {
	Renderable
	GetRGBA() *image.RGBA
	Modify(...Modification) Modifiable
	Copy() Modifiable
}

// And chains together multiple Modifications into a single Modification
func And(ms ...Modification) Modification {
	return func(rgba image.Image) *image.RGBA {
		rgba2 := ms[0](rgba)
		for i := 1; i < len(ms); i++ {
			rgba2 = ms[i](rgba2)
		}
		return rgba2
	}
}

// Brighten brightens an image
func Brighten(brightenBy float32) Modification {
	return func(rgba image.Image) *image.RGBA {
		filter := gift.New(
			gift.Brightness(brightenBy))
		dst := image.NewRGBA(filter.Bounds(rgba.Bounds()))
		filter.Draw(dst, rgba)
		return dst
	}
}

// FlipX returns a new rgba which is flipped
// over the horizontal axis.
func FlipX(rgba image.Image) *image.RGBA {
	filter := gift.New(
		gift.FlipHorizontal())
	dst := image.NewRGBA(filter.Bounds(rgba.Bounds()))
	filter.Draw(dst, rgba)
	return dst
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

// Fade reduces the alpha of an image
func Fade(alpha int) Modification {
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
func ApplyColor(c color.Color) Modification {
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

// ColorBalance takes in 3 numbers between -100 and 500 and applies it to the given image
func ColorBalance(r, g, b float32) Modification {
	return func(rgba image.Image) *image.RGBA {
		filter := gift.New(gift.ColorBalance(r, g, b))
		dst := image.NewRGBA(filter.Bounds(rgba.Bounds()))
		filter.Draw(dst, rgba)
		return dst
	}
}

// FillMask replaces alpha 0 pixels in an RGBA with corresponding
// pixels in a second RGBA.
func FillMask(img image.RGBA) Modification {
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
func ApplyMask(img image.RGBA) Modification {
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

// Rotate returns a rotated rgba.
func Rotate(degrees int) Modification {
	return func(rgba image.Image) *image.RGBA {
		filter := gift.New(
			gift.Rotate(float32(degrees), transparent, gift.CubicInterpolation))
		dst := image.NewRGBA(filter.Bounds(rgba.Bounds()))
		filter.Draw(dst, rgba)
		return dst
	}
}

// Scale returns a scaled rgba.
func Scale(xRatio, yRatio float64) Modification {
	return func(rgba image.Image) *image.RGBA {
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
}
