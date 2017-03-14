package render

import (
	"bitbucket.org/oakmoundstudio/oak/physics"
	"image"
	"image/color"
	"math"
)

type GradientBox struct {
	Sprite
}

var (
	HorizontalProgress = func(x, y, w, h int) float64 {
		return float64(x) / float64(w)
	}
	VerticalProgress = func(x, y, w, h int) float64 {
		return float64(y) / float64(h)
	}

	CircularProgress = func(x, y, w, h int) float64 {
		xRadius := float64(w) / 2
		yRadius := float64(h) / 2
		dX := math.Abs(float64(x) - xRadius)
		dY := math.Abs(float64(y) - yRadius)
		progress := math.Pow(dX/xRadius, 2) + math.Pow(dY/yRadius, 2)
		if progress > 1 {
			progress = 1
		}
		return progress
	}
)

func NewGradientBox(w, h int, startColor, endColor color.Color, pFunction progressFunction) *GradientBox {
	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)

	r, g, b, a := startColor.RGBA()
	r2, g2, b2, a2 := endColor.RGBA()

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			progress := pFunction(x, y, w, h)
			c := color.RGBA64{
				uint16OnScale(r, r2, progress),
				uint16OnScale(g, g2, progress),
				uint16OnScale(b, b2, progress),
				uint16OnScale(a, a2, progress),
			}
			rgba.Set(x, y, c)
		}
	}
	return &GradientBox{
		Sprite{
			LayeredPoint: LayeredPoint{
				Vector: physics.Vector{
					X: 0.0,
					Y: 0.0,
				},
			},
			r: rgba,
		},
	}
}

func NewHorizontalGradientBox(w, h int, startColor, endColor color.Color) *GradientBox {
	return NewGradientBox(w, h, startColor, endColor, HorizontalProgress)
}

func NewVerticalGradientBox(w, h int, startColor, endColor color.Color) *GradientBox {
	return NewGradientBox(w, h, startColor, endColor, VerticalProgress)
}

func NewCircularGradientBox(w, h int, startColor, endColor color.Color) *GradientBox {
	return NewGradientBox(w, h, startColor, endColor, CircularProgress)
}

func uint16OnScale(n, endN uint32, progress float64) uint16 {
	return uint16((float64(endN) - float64(endN)*(1.0-progress) + float64(n)*(1.0-progress)))
}

type progressFunction func(x, y, w, h int) float64
