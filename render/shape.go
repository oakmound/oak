package render

import (
	"image"
	"image/color"

	"github.com/oakmound/oak/v3/shape"
)

// SpriteFromShape converts a shape into a sprite with values in the shape being colored
// 'on' and out of the shape being colored 'off'
func SpriteFromShape(sh shape.Shape, w, h int, on, off color.Color) *Sprite {
	rect := sh.Rect(w, h)
	rgba := image.NewRGBA(image.Rect(0, 0, len(rect), len(rect[0])))
	sp := NewSprite(0, 0, rgba)
	for x := 0; x < len(rect); x++ {
		for y := 0; y < len(rect[0]); y++ {
			if rect[x][y] {
				sp.Set(x, y, on)
			} else {
				sp.Set(x, y, off)
			}
		}
	}
	return sp
}
