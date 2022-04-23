//go:build !nogift
// +build !nogift

package scene

import (
	"image"

	"github.com/oakmound/oak/v4/render/mod"
)

// Fade is a scene transition that fades to black at a given rate for
// a total of 'frames' frames
func Fade(rate float32, frames int) Transition {
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
