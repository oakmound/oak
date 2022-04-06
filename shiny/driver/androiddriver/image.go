//go:build android
// +build android

package androiddriver

import (
	"image"
	"sync"

	"golang.org/x/mobile/exp/gl/glutil"
)

type imageImpl struct {
	screen   *Screen
	size     image.Point
	img      *glutil.Image
	deadLock sync.Mutex
	dead     bool
}

func (ii *imageImpl) Size() image.Point {
	return ii.size
}

func (ii *imageImpl) Bounds() image.Rectangle {
	return image.Rect(0, 0, ii.size.X, ii.size.Y)
}

func (ii *imageImpl) Release() {
	ii.deadLock.Lock()
	ii.img.Release()
	ii.dead = true
	ii.deadLock.Unlock()
}

func (ii *imageImpl) RGBA() *image.RGBA {
	ii.deadLock.Lock()
	if ii.dead {
		ii.img = ii.screen.images.NewImage(ii.size.X, ii.size.Y)
		ii.dead = false
	}
	ii.deadLock.Unlock()
	return ii.img.RGBA
}
