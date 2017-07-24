package render

import (
	"image"
	"image/color"
	"math"
)

// NewLine returns a line from x1,y1 to x2,y2 with the given color
func NewLine(x1, y1, x2, y2 float64, c color.Color) *Sprite {
	return NewThickLine(x1, y1, x2, y2, c, 0)
}

// NewThickLine returns a Line that has some value of thickness
func NewThickLine(x1, y1, x2, y2 float64, c color.Color, thickness int) *Sprite {

	var rgba *image.RGBA
	// We subtract the minimum from each side here
	// to normalize the new line segment toward the origin
	minX := math.Min(x1, x2)
	minY := math.Min(y1, y2)
	rgba = drawLineBetween(int(x1-minX), int(y1-minY), int(x2-minX), int(y2-minY), c, thickness)

	return NewSprite(minX-float64(thickness), minY-float64(thickness), rgba)
}

// DrawLineOnto draws a line onto an image rgba from one point to another
// Todo: this and drawLineBetween should be combined to reduce duplicate code
func DrawLineOnto(rgba *image.RGBA, x1, y1, x2, y2 int, c color.Color) {

	xDelta := math.Abs(float64(x2 - x1))
	yDelta := math.Abs(float64(y2 - y1))

	xSlope := -1
	if x2 < x1 {
		xSlope = 1
	}
	ySlope := -1
	if y2 < y1 {
		ySlope = 1
	}

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
}

func drawLineBetween(x1, y1, x2, y2 int, c color.Color, th int) *image.RGBA {

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
	// Todo: document why we add one here
	h := int(yDelta) + 1
	if y2 < y1 {
		ySlope = 1
	}
	rect := image.Rect(0, 0, int(xDelta)+2*th, h+2*th)
	rgba := image.NewRGBA(rect)

	x2 += th
	y2 += th
	x1 += th
	y1 += th

	err := xDelta - yDelta
	var err2 float64
	for i := 0; true; i++ {

		for xm := x2 - th; xm <= (x2 + th); xm++ {
			for ym := y2 - th; ym <= (y2 + th); ym++ {
				rgba.Set(xm, ym, c)
			}
		}
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
