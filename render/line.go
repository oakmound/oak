package render

import (
	"bitbucket.org/oakmoundstudio/oak/physics"
	"image"
	"image/color"
	"math"
)

type Line struct {
	Sprite
}

func NewLine(x1, y1, x2, y2 float64, c color.Color) *Line {

	var rgba *image.RGBA
	// We subtract the minimum from each side here
	// to normalize the new line segment toward the origin
	minX := math.Min(x1, x2)
	minY := math.Min(y1, y2)
	rgba = drawLineBetween(int(x1-minX), int(y1-minY), int(x2-minX), int(y2-minY), c)

	return &Line{
		Sprite{
			LayeredPoint: LayeredPoint{
				Vector: physics.Vector{
					X: minX,
					Y: minY,
				},
			},
			r: rgba,
		},
	}
}

func drawLineBetween(x1, y1, x2, y2 int, c color.Color) *image.RGBA {

	// Bresenham's line-drawing algorithm from wikipedia
	xDelta := math.Abs(float64(x2 - x1))
	yDelta := math.Abs(float64(y2 - y1))

	if xDelta == 0 && yDelta == 0 {
		rect := image.Rect(0, 0, 1, 1)
		rgba := image.NewRGBA(rect)
		rgba.Set(0, 0, c)
		return rgba
	} else if xDelta == 0 {
		rect := image.Rect(0, 0, 1, int(math.Floor(yDelta)))
		rgba := image.NewRGBA(rect)
		for i := 0; i < int(math.Floor(yDelta)); i++ {
			rgba.Set(0, i, c)
		}
		return rgba
	}

	xSlope := -1
	if x2 < x1 {
		xSlope = 1
	}
	ySlope := -1
	h := int(yDelta) + 1
	if y2 < y1 {
		ySlope = 1
	}
	rect := image.Rect(0, 0, int(xDelta), h)
	rgba := image.NewRGBA(rect)

	err := xDelta - yDelta
	var err2 float64
	for i := 0; true; i++ {

		rgba.Set(x2, y2, c)
		if x2 == x1 && y2 == y1 {
			break
		}
		err2 = 2 * err
		if err2 > -1*yDelta {
			err -= yDelta
			x2 += xSlope
		}
		if err2 < xDelta {
			err += xDelta
			y2 += ySlope
		}
	}

	return rgba
}
