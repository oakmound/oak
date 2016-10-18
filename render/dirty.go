package render

import (
	"fmt"
)

const (
	dirtyWidth  = 64
	dirtyHeight = 64
	// World size 4000
	dirtyZonesX = 4096 / dirtyWidth
	dirtyZonesY = 4096 / dirtyHeight
)

var (
	dirtyZones = [dirtyZonesX][dirtyZonesY]bool{}
)

func dirtyBounds(x, y int) (int, int) {
	x /= dirtyWidth
	y /= dirtyHeight
	if x <= 0 {
		x = 1
	} else if x >= dirtyZonesX {
		x = dirtyZonesX - 1
	}
	if y <= 0 {
		y = 1
	} else if y >= dirtyZonesY {
		y = dirtyZonesY - 1
	}
	return x, y
}

func SetDirty(fx, fy float64) {
	x := int(fx)
	y := int(fy)
	x, y = dirtyBounds(x, y)
	fmt.Println(x, y)
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
