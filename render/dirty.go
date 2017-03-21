package render

import (
	"image"
)

const (
	DirtyWidth  = 64
	DirtyHeight = 64
	// World size 4000
	DirtyZonesX = 4096 / DirtyWidth
	DirtyZonesY = 4096 / DirtyHeight
)

var (
	DirtyBounds = image.Rect(0, 0, DirtyWidth, DirtyHeight)
	dirtyZones  = [DirtyZonesX][DirtyZonesY]bool{}
)

func dirtyBounds(x, y int) (int, int) {
	x /= DirtyWidth
	y /= DirtyHeight
	if x <= 0 {
		x = 1
	} else if x >= DirtyZonesX {
		x = DirtyZonesX - 1
	}
	if y <= 0 {
		y = 1
	} else if y >= DirtyZonesY {
		y = DirtyZonesY - 1
	}
	return x, y
}

func SetDirty(fx, fy float64) {
	x := int(fx)
	y := int(fy)
	x, y = dirtyBounds(x, y)
	for i := x - 1; i < x+1; i++ {
		for j := y - 1; j < y+1; j++ {
			dirtyZones[i][j] = true
		}
	}
}

func IsDirty(x, y int) bool {
	x, y = dirtyBounds(x, y)
	for i := x - 1; i < x+1; i++ {
		for j := y - 1; j < y+1; j++ {
			if dirtyZones[i][j] {
				return true
			}
		}
	}
	return false
}
