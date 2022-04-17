package scene

import (
	"image"
	"image/draw"

	"github.com/oakmound/oak/v3/render/mod"
)

// Transition functions can be set to occur at the end of a scene.
type Transition func(*image.RGBA, int) bool

// Zoom transitions by performing a simplistic zoom each frame towards some
// percentage-based part of the screen.
func Zoom(xPerc, yPerc float64, frames int, zoomRate float64) Transition {
	return func(buf *image.RGBA, frame int) bool {
		if frame > frames {
			return false
		}
		z := mod.Zoom(xPerc, yPerc, 1+zoomRate*float64(frame))
		draw.Draw(buf, buf.Bounds(), z(buf), image.ZP, draw.Src)
		return true
	}
}
