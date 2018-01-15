package shape

import (
	"errors"
	"math"

	"github.com/oakmound/oak/alg/intgeom"
)

const (
	top = iota
	topright
	right
	bottomright
	bottom
	bottomleft
	left
	topleft
	lastdirection
)

var (
	xyMods = []int{
		0, -1,
		1, -1,
		1, 0,
		1, 1,
		0, 1,
		-1, 1,
		-1, 0,
		-1, -1,
	}
	pointDeltas = []int{
		1, 0,
		0, 1,
		0, 1,
		-1, 0,
		-1, 0,
		0, -1,
		0, -1,
		1, 0,
	}
)

// ToOutline returns the set of points along the input shape's outline, if
// one exists.
func ToOutline(shape Shape) func(...int) ([]intgeom.Point, error) {
	return func(sizes ...int) ([]intgeom.Point, error) {
		return toOutline(shape, 1, sizes...)
	}
}

func parseSizes(sizes []int) (int, int) {
	w := sizes[0]
	h := sizes[0]
	if len(sizes) > 1 {
		h = sizes[1]
	}
	return w, h
}

// this is a hack to support 4 and 8 directional outlines
func toOutline(shape Shape, dirInc int, sizes ...int) ([]intgeom.Point, error) {
	w, h := parseSizes(sizes)

	//First decrement on diagonal to find start of outline
	startX := 0.0
	startY := 0.0
	fw := float64(w)
	fh := float64(h)
	maxDim := math.Max(fw, fh)
	xDelta := fw / maxDim
	yDelta := fh / maxDim
	for !shape.In(int(startX), int(startY), sizes...) {
		startX += xDelta
		startY += yDelta
		if startX >= fw || startY >= fh {
			return []intgeom.Point{}, errors.New("Could not find an outline space on the shape's diagonal")
		}
	}

	for startY >= 0 && shape.In(int(startX), int(startY), sizes...) {
		startY--
	}
	startY++

	//Here we have found a point on the outline
	sx := int(startX)
	sy := int(startY)
	x := sx
	y := sy

	outline := []intgeom.Point{intgeom.NewPoint(x, y)}

	direction := topright
	for i := 1; i < dirInc; i++ {
		direction++
	}

	x += xyMods[direction*2]
	y += xyMods[direction*2+1]

	for direction != top && !inOutline(shape, x, y, w, h) {
		for i := 0; i < dirInc; i++ {
			x += pointDeltas[direction*2]
			y += pointDeltas[direction*2+1]
			direction = (direction + 1) % lastdirection
		}
	}
	if direction == top {
		return outline, nil
	}

	return followOutline(shape, dirInc, x, y, sx, sy, w, h, direction, outline), nil
}

func followOutline(shape Shape, dirInc int, x, y, sx, sy, w, h, direction int, outline []intgeom.Point) []intgeom.Point {
	//Follow the outline point by point
	for x != sx || y != sy {
		outline = append(outline, intgeom.NewPoint(x, y))
		direction -= 2
		if direction < 0 {
			direction += lastdirection
		}
		x += xyMods[direction*2]
		y += xyMods[direction*2+1]
		//From a point on the outline look clockwise around for next direction
		for !inOutline(shape, x, y, w, h) {
			for i := 0; i < dirInc; i++ {
				x += pointDeltas[direction*2]
				y += pointDeltas[direction*2+1]
				direction = (direction + 1) % lastdirection
			}
		}
	}
	return outline
}

// ToOutline4 returns the set of points along the input shape's outline, if
// one exists, but will move only up, left, right, or down to form this outline.
func ToOutline4(shape Shape) func(...int) ([]intgeom.Point, error) {
	return func(sizes ...int) ([]intgeom.Point, error) {
		return toOutline(shape, 2, sizes...)
	}
}

func inOutline(s Shape, x, y, w, h int) bool {
	return (x < w && x >= 0 && y < h && y >= 0) && s.In(x, y, w, h)
}
