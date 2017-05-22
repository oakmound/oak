package render

import (
	"image"
	"image/color"
	"image/draw"
)

// NewColorBox returns a Sprite full of a given color with the given dimensions
func NewColorBox(w, h int, c color.Color) *Sprite {
	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, image.NewUniform(c), image.Point{0, 0}, draw.Src)
	return NewSprite(0, 0, rgba)
}
