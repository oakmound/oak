package render

import (
	simplex "github.com/ojrac/opensimplex-go"
	"image"
	"image/color"
	// "image/draw"
	"time"
)

type NoiseBox struct {
	Sprite
}

func NewNoiseBox(w, h int) *NoiseBox {
	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)
	noise := simplex.NewWithSeed(time.Now().Unix())

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			scale := uint8(noise.Eval2(float64(x), float64(y)) * 255)
			rgba.Set(x, y, color.RGBA{scale, scale, scale, 255})
		}
	}

	return &NoiseBox{
		Sprite{
			LayeredPoint: LayeredPoint{
				Point: Point{
					X: 0.0,
					Y: 0.0,
				},
			},
			r: rgba,
		},
	}
}
