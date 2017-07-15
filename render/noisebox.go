package render

import (
	"image"
	"image/color"
	"math/rand"
	"time"
)

// NewNoiseBox returns a box of noise
func NewNoiseBox(w, h int) *Sprite {
	return NewSeededNoiseBox(w, h, time.Now().Unix())
}

// NewSeededNoiseBox returns a box of noise seeded at a specific value
// this previously used a complex noise function, but this refused to
// run on windows 32bit and was overkill, so it now uses math/rand
func NewSeededNoiseBox(w, h int, seed int64) *Sprite {
	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)
	rng := rand.New(rand.NewSource(seed))

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			v := uint8(rng.Intn(256))
			rgba.Set(x, y, color.RGBA{v, v, v, 255})
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
