package render

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/oakmound/oak/v4/alg/intgeom"
)

// NewColorBox returns a Sprite full of a given color with the given dimensions
// Deprecated: Use NewColorboxM (for a Modifiable) or NewColorBoxR.
func NewColorBox(w, h int, c color.Color) *Sprite {
	return NewColorBoxM(w, h, c)
}

// NewColorBoxM returns a modifiable Color Box (as a Sprite)
func NewColorBoxM(w, h int, c color.Color) *Sprite {
	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, image.NewUniform(c), image.Point{0, 0}, draw.Src)
	return NewSprite(0, 0, rgba)
}

// ColorBoxR is a renderable color box. It is a smaller structure and should
// render faster than a ColorBoxM.
type ColorBoxR struct {
	LayeredPoint
	Dims  intgeom.Point2
	Color *image.Uniform
}

// NewColorBoxR creates a color box. Colorboxes made without using this constructor
// may not function.
func NewColorBoxR(w, h int, c color.Color) *ColorBoxR {
	return &ColorBoxR{
		LayeredPoint: NewLayeredPoint(0, 0, 0),
		Dims:         intgeom.Point2{w, h},
		Color:        image.NewUniform(c),
	}
}

// GetDims returns the dimensiosn of this colorbox
func (cb *ColorBoxR) GetDims() (int, int) {
	return cb.Dims.X(), cb.Dims.Y()
}

// Draw renders this colorbox to screen.
func (cb *ColorBoxR) Draw(buff draw.Image, xOff, yOff float64) {
	pt := image.Point{int((cb.X() + xOff)), int((cb.Y() + yOff))}
	max := pt.Add(image.Point{X: cb.Dims.X(), Y: cb.Dims.Y()})
	draw.Draw(
		buff,
		image.Rectangle{Min: pt, Max: max},
		cb.Color,
		image.Point{0, 0},
		draw.Over)
}
