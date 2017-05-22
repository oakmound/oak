// OpenSimplex-Go refuses to compile for Windows,386

// +build !386

package render

import (
	"image"
	"image/color"
	"time"

	simplex "github.com/ojrac/opensimplex-go"
)

// NewNoiseBox returns a box of noise
func NewNoiseBox(w, h int) *Sprite {
	return NewSeededNoiseBox(w, h, time.Now().Unix())
}

// NewSeededNoiseBox returns a box of noise seeded at a specific value
func NewSeededNoiseBox(w, h int, seed int64) *Sprite {
	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)
	noise := simplex.NewWithSeed(seed)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			scale := uint8(noise.Eval2(float64(x), float64(y)) * 255)
			rgba.Set(x, y, color.RGBA{scale, scale, scale, 255})
		}
	}

	return NewSprite(0, 0, rgba)
}

// NewNoiseSequence returns a sequence of noise boxes
func NewNoiseSequence(w, h, frames int, fps float64) *Sequence {
	mods := make([]Modifiable, frames)
	for i := 0; i < frames; i++ {
		mods[i] = NewSeededNoiseBox(w, h, time.Now().Unix()*int64(i))
	}
	return NewSequence(mods, fps)
}
