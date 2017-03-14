// OpenSimplex-Go refuses to compile for Windows,386

// +build !386

package render

import (
	"image"
	"image/color"
	"time"

	"bitbucket.org/oakmoundstudio/oak/physics"
	simplex "github.com/ojrac/opensimplex-go"
)

type NoiseBox struct {
	Sprite
}

func NewNoiseBox(w, h int) *NoiseBox {
	return NewSeededNoiseBox(w, h, time.Now().Unix())
}

func NewSeededNoiseBox(w, h int, seed int64) *NoiseBox {
	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)
	noise := simplex.NewWithSeed(seed)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			scale := uint8(noise.Eval2(float64(x), float64(y)) * 255)
			rgba.Set(x, y, color.RGBA{scale, scale, scale, 255})
		}
	}

	return &NoiseBox{
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

func NewNoiseSequence(w, h, frames int, fps float64) *Sequence {
	mods := make([]Modifiable, frames)
	for i := 0; i < frames; i++ {
		mods[i] = NewSeededNoiseBox(w, h, time.Now().Unix()*int64(i))
	}
	return NewSequence(mods, fps)
}
