package render

import (
	// This file is being slowly converted to use gift over manual math and loops,
	// because our math / loops will be more likely to have (and have already had)
	// missable bugs.
	//"github.com/anthonynsimon/bild/blend"

	"github.com/oakmound/oak/alg"

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

// todo: this should not be in this package
type point struct {
	X, Y float64
}

// CutRound rounds the edges of the Modifiable with Bezier curves.
// Todo: A nice bezier curve toolset would be nice
func CutRound(xOff, yOff float64) Modification {
	return func(rgba image.Image) *image.RGBA {
		bds := rgba.Bounds()
		newRgba := image.NewRGBA(bds)

		// start off as a copy
		for x := bds.Min.X; x < bds.Max.X; x++ {
			for y := bds.Min.Y; y < bds.Max.Y; y++ {
				newRgba.Set(x, y, rgba.At(x, y))
			}
		}
		// For each corner, define directions

		corners := [][4]int{
			// X, Y, xDir, yDir
			{bds.Min.X, bds.Min.Y, 1, 1},
			{bds.Min.X, bds.Max.Y, 1, -1},
			{bds.Max.X, bds.Max.Y, -1, -1},
			{bds.Max.X, bds.Min.Y, -1, 1},
		}
		for _, c := range corners {
			// 3 point Bezier curve
			x1 := float64(c[0])
			y1 := float64(c[1])
			x2 := x1 + (float64(bds.Max.X*c[2]) * xOff)
			y2 := y1 + (float64(bds.Max.Y*c[3]) * yOff)
			p1 := point{x2, y1}
			p2 := point{x1, y1}
			p3 := point{x1, y2}
			//fmt.Println("Corners", p1, p2, p3)

			// Progressing along the curve, whenever a new y value is
			// intersected at a pixel delete all values
			// from the image above(or below, for negative c[3])
			// that pixel

			// todo: non-arbitrary progress increment
			for progress := 0.0; progress < 1.0; progress += 0.01 {
				p4 := pointBetween(p1, p2, progress)
				p5 := pointBetween(p2, p3, progress)
				curveAt := pointBetween(p4, p5, progress)
				//fmt.Println("Curve, progress:", progress, "pts", p4, p5, curveAt)

				// Could only redo this loop at new y values to save time,
				// but because this is currently just a pre-processing modification
				// it should be okay
				x := alg.RoundF64(curveAt.X)
				for y := alg.RoundF64(curveAt.Y); y <= bds.Max.Y && y >= bds.Min.Y; y -= c[3] {
					newRgba.Set(x, y, color.RGBA{0, 0, 0, 0})
				}
			}
		}

		return newRgba
	}
}

// todo: this should not be in this package
func pointBetween(p1, p2 point, f float64) point {
	return point{p1.X*(1-f) + p2.X*f, p1.Y*(1-f) + p2.Y*f}
}

// CutRel acts like Cut, but takes in a multiplier on the
// existing dimensions of the image.
func CutRel(relWidth, relHeight float64) Modification {
	return func(rgba image.Image) *image.RGBA {
		bds := rgba.Bounds()
		newWidth := alg.RoundF64(float64(bds.Max.X) * relWidth)
		newHeight := alg.RoundF64(float64(bds.Max.Y) * relHeight)
		newRgba := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		for x := 0; x < newWidth; x++ {
			for y := 0; y < newHeight; y++ {
				newRgba.Set(x, y, rgba.At(x, y))
			}
		}
		return newRgba
	}
}

// Cut reduces (or increases, adding nothing)
// the dimensions of the input image, setting them to newWidth and
// newHeight. (Consider: use generic int modifiers here so we
// don't need CutRel and Cut? i.e a function header like
// Cut(wMod, hMod func(int) int)? )
func Cut(newWidth, newHeight int) Modification {
	return func(rgba image.Image) *image.RGBA {
		newRgba := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		for x := 0; x < newWidth; x++ {
			for y := 0; y < newHeight; y++ {
				newRgba.Set(x, y, rgba.At(x, y))
			}
		}
		return newRgba
	}
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
		filter := gift.New(
			gift.Resize(w, h, gift.CubicResampling))
		dst := image.NewRGBA(filter.Bounds(rgba.Bounds()))
		filter.Draw(dst, rgba)
		return dst
	}
}
