package scene

import (
	"image"
	"image/draw"

	"github.com/oakmound/oak/render/mod"
)

var (
	zeroPoint = image.Point{X: 0, Y: 0}
)

// Transition functions can be set to occur at the end of a scene.
type Transition func(*image.RGBA, int) bool

// Fade is a scene transition that fades to black at a given rate for
// a total of 'frames' frames
func Fade(rate float32, frames int) func(*image.RGBA, int) bool {
	rate *= -1
	return func(buf *image.RGBA, frame int) bool {
		if frame > frames {
			return false
		}
		i := float32(frame)
		mod.Brighten(rate * i)(buf)
		return true
	}
}

// Zoom transitions by performing a simplistic zoom each frame towards some
// percentange-based part of the screen.
func Zoom(xPerc, yPerc float64, frames int, zoomRate float64) func(*image.RGBA, int) bool {
	return func(buf *image.RGBA, frame int) bool {
		if frame > frames {
			return false
		}
		z := mod.Zoom(xPerc, yPerc, 1+zoomRate*float64(frame))
		draw.Draw(buf, buf.Bounds(), z(buf), zeroPoint, draw.Src)
		return true
	}
}
