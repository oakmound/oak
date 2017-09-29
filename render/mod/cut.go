package mod

import (
	"image"
	"image/color"

	"github.com/oakmound/oak/alg"
	"github.com/oakmound/oak/alg/floatgeom"
)

// CutRound rounds the edges of the Modifiable with Bezier curves.
// Todo: We have a nice bezier toolkit now, so use it here
func CutRound(xOff, yOff float64) Mod {
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
			p1 := floatgeom.Point2{x2, y1}
			p2 := floatgeom.Point2{x1, y1}
			p3 := floatgeom.Point2{x1, y2}

			// Progressing along the curve, whenever a new y value is
			// intersected at a pixel delete all values
			// from the image above(or below, for negative c[3])
			// that pixel

			// todo: non-arbitrary progress increment
			for progress := 0.0; progress < 1.0; progress += 0.01 {
				p4 := pointBetween(p1, p2, progress)
				p5 := pointBetween(p2, p3, progress)
				curveAt := pointBetween(p4, p5, progress)

				// Could only redo this loop at new y values to save time,
				// but because this is currently just a pre-processing Mod
				// it should be okay
				x := alg.RoundF64(curveAt.X())
				for y := alg.RoundF64(curveAt.Y()); y <= bds.Max.Y && y >= bds.Min.Y; y -= c[3] {
					newRgba.Set(x, y, color.RGBA{0, 0, 0, 0})
				}
			}
		}

		return newRgba
	}
}

// todo: this should not be in this package
func pointBetween(p1, p2 floatgeom.Point2, f float64) floatgeom.Point2 {
	return floatgeom.Point2{p1.X()*(1-f) + p2.X()*f, p1.Y()*(1-f) + p2.Y()*f}
}

//CutFn  can reduce or add blank space to an input image.
// Each input function decides the starting location or offset of a cut.
func CutFn(xMod, yMod, wMod, hMod func(int) int) Mod {
	return func(rgba image.Image) *image.RGBA {
		bds := rgba.Bounds()
		startX := xMod(bds.Max.X)
		startY := yMod(bds.Max.Y)
		newWidth := wMod(bds.Max.X)
		newHeight := hMod(bds.Max.Y)

		newRgba := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		for x := 0; x < newWidth; x++ {
			for y := 0; y < newHeight; y++ {
				newRgba.Set(x, y, rgba.At(x+startX, y+startY))
			}
		}
		return newRgba
	}
}

// CutFromLeft acts like cut but removes from the left hand side rather than the right
func CutFromLeft(newWidth, newHeight int) Mod {
	return CutFn(func(w int) int {
		out := w - newWidth
		return out
	},
		func(h int) int {
			out := h - newHeight
			return out
		},
		func(int) int {
			return newWidth
		},
		func(int) int {
			return newHeight
		})
}

// CutRel acts like Cut, but takes in a multiplier on the
// existing dimensions of the image.
func CutRel(relWidth, relHeight float64) Mod {
	return CutFn(func(int) int { return 0 },
		func(int) int { return 0 },
		func(w int) int { return alg.RoundF64(float64(w) * relWidth) },
		func(h int) int { return alg.RoundF64(float64(h) * relHeight) })
}

// Cut reduces (or increases, adding nothing)
// the dimensions of the input image, setting them to newWidth and
// newHeight. (Consider: use generic int modifiers here so we
// don't need CutRel and Cut? i.e a function header like
// Cut(wMod, hMod func(int) int)? )
func Cut(newWidth, newHeight int) Mod {
	return CutFn(func(int) int { return 0 },
		func(int) int { return 0 },
		func(int) int { return newWidth },
		func(int) int { return newHeight })
}
