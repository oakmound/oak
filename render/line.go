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
	rgba = drawLineBetween(int(x1-minX), int(y1-minY), int(x2-minX), int(y2-minY), func(rgba *image.RGBA, totalD float64, progress, x, y int) { rgba.Set(x, y, c) }, thickness)

	return NewSprite(minX-float64(thickness), minY-float64(thickness), rgba)
}

// NewGradientLine returns a Line that has some value of thickness along with a start and end color
func NewGradientLine(x1, y1, x2, y2 float64, c1, c2 color.Color, thickness int) *Sprite {

	var rgba *image.RGBA
	// We subtract the minimum from each side here
	// to normalize the new line segment toward the origin
	minX := math.Min(x1, x2)
	minY := math.Min(y1, y2)
	colorer := func(rgba *image.RGBA, totalD float64, progress, x, y int) {
		percentProgress := float64(progress) / totalD
		c := GradientColorAt(c1, c2, percentProgress)
		rgba.Set(x, y, c)
	}
	rgba = drawLineBetween(int(x1-minX), int(y1-minY), int(x2-minX), int(y2-minY), colorer, thickness)

	return NewSprite(minX-float64(thickness), minY-float64(thickness), rgba)
}

// DrawThickLine acts like DrawlineOnto, but takes in thickness of the given line
func DrawThickLine(rgba *image.RGBA, x1, y1, x2, y2 int, c color.Color, thickness int) {
	drawLine(rgba, x1, y1, x2, y2, thickness, func(rgba *image.RGBA, totalD float64, progress, x, y int) { rgba.Set(x, y, c) })
}

//DrawGradientLine acts like DrawThickLine but also applies a gradent to the line
func DrawGradientLine(rgba *image.RGBA, x1, y1, x2, y2 int, c1, c2 color.Color, thickness int) {
	colorer := func(rgba *image.RGBA, totalD float64, progress, x, y int) {
		percentProgress := float64(progress) / totalD
		c := GradientColorAt(c1, c2, percentProgress)
		rgba.Set(x, y, c)
	}
	drawLine(rgba, x1, y1, x2, y2, thickness, colorer)
}

// DrawLineOnto draws a line onto an image rgba from one point to another
// Todo: this and drawLineBetween should be combined to reduce duplicate code
func DrawLineOnto(rgba *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
	DrawThickLine(rgba, x1, y1, x2, y2, c, 0)
}

type pixelColorer func(rgba *image.RGBA, totalDistance float64, progress, x, y int)

func drawLine(rgba *image.RGBA, x1, y1, x2, y2 int, thickness int, colorer pixelColorer) {
	xDelta := math.Abs(float64(x2 - x1))
	yDelta := math.Abs(float64(y2 - y1))
	totalDelta := math.Sqrt(xDelta*xDelta + yDelta + yDelta)
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
		for xm := x2 - thickness; xm <= (x2 + thickness); xm++ {
			for ym := y2 - thickness; ym <= (y2 + thickness); ym++ {
				colorer(rgba, totalDelta, i, xm, ym)
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
}

func drawLineBetween(x1, y1, x2, y2 int, colorer pixelColorer, th int) *image.RGBA {

	// Bresenham's line-drawing algorithm from wikipedia
	xDelta := math.Abs(float64(x2 - x1))
	yDelta := math.Abs(float64(y2 - y1))
	totalDelta := math.Sqrt(xDelta*xDelta + yDelta + yDelta)

	if xDelta == 0 && yDelta == 0 {
		rect := image.Rect(0, 0, 1, 1)
		rgba := image.NewRGBA(rect)
		colorer(rgba, totalDelta, 0, 0, 0)
		return rgba
	} else if xDelta == 0 {
		rect := image.Rect(0, 0, 1, int(math.Floor(yDelta)))
		rgba := image.NewRGBA(rect)
		for i := 0; i < int(math.Floor(yDelta)); i++ {
			colorer(rgba, totalDelta, i, 0, i)
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
				colorer(rgba, totalDelta, i, xm, ym)

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
