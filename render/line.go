package render

import (
	"image"
	"image/color"
	"math"

	"github.com/200sc/go-dist/colorrange"
)

// Todo:
// Line and color optional syntax? should color / progress move to separate packages?

// NewLine returns a line from x1,y1 to x2,y2 with the given color
func NewLine(x1, y1, x2, y2 float64, c color.Color) *Sprite {
	return NewThickLine(x1, y1, x2, y2, c, 0)
}

// NewThickLine returns a Line that has some value of thickness
func NewThickLine(x1, y1, x2, y2 float64, c color.Color, thickness int) *Sprite {
	return NewLineColored(x1, y1, x2, y2, IdentityColorer(c), thickness)
}

// NewGradientLine returns a Line that has some value of thickness along with a start and end color
func NewGradientLine(x1, y1, x2, y2 float64, c1, c2 color.Color, thickness int) *Sprite {
	colorer := colorrange.NewLinear(c1, c2).Percentile
	return NewLineColored(x1, y1, x2, y2, colorer, thickness)
}

// NewLineColored returns a line with a custom function for how each pixel in that line should be colored.
func NewLineColored(x1, y1, x2, y2 float64, colorer Colorer, thickness int) *Sprite {
	var rgba *image.RGBA
	// We subtract the minimum from each side here
	// to normalize the new line segment toward the origin
	minX := math.Min(x1, x2)
	minY := math.Min(y1, y2)
	rgba = drawLineBetween(int(x1-minX), int(y1-minY), int(x2-minX), int(y2-minY), colorer, thickness)
	return NewSprite(minX-float64(thickness), minY-float64(thickness), rgba)
}

// DrawLineOnto draws a line onto an image rgba from one point to another
// Todo: 2.0, rename to DrawLine
func DrawLineOnto(rgba *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
	DrawThickLine(rgba, x1, y1, x2, y2, c, 0)
}

// DrawThickLine acts like DrawlineOnto, but takes in thickness of the given line
func DrawThickLine(rgba *image.RGBA, x1, y1, x2, y2 int, c color.Color, thickness int) {
	DrawLineColored(rgba, x1, y1, x2, y2, thickness, IdentityColorer(c))
}

//DrawGradientLine acts like DrawThickLine but also applies a gradient to the line
func DrawGradientLine(rgba *image.RGBA, x1, y1, x2, y2 int, c1, c2 color.Color, thickness int) {
	colorer := colorrange.NewLinear(c1, c2).Percentile
	DrawLineColored(rgba, x1, y1, x2, y2, thickness, colorer)
}

// DrawLineColored acts like DrawThickLine, but takes in a custom colorer function for how it draws its line.
func DrawLineColored(rgba *image.RGBA, x1, y1, x2, y2, thickness int, colorer Colorer) {

	xDelta := math.Abs(float64(x2 - x1))
	yDelta := math.Abs(float64(y2 - y1))

	xSlope := -1
	x3 := x1
	if x2 < x1 {
		xSlope = 1
		x3 = x2
	}
	ySlope := -1
	y3 := y1
	if y2 < y1 {
		ySlope = 1
		y3 = y2
	}

	w := int(xDelta)
	h := int(yDelta)

	progress := func(x, y, w, h int) float64 {
		wprg := HorizontalProgress(x, y, w, h)
		vprg := VerticalProgress(x, y, w, h)
		if ySlope == -1 {
			vprg = 1 - vprg
		}
		if xSlope == -1 {
			wprg = 1 - wprg
		}
		return wprg + vprg/2
	}

	err := xDelta - yDelta
	var err2 float64
	for i := 0; true; i++ {
		for xm := x2 - thickness; xm <= (x2 + thickness); xm++ {
			for ym := y2 - thickness; ym <= (y2 + thickness); ym++ {
				p := progress(xm-x3, ym-y3, w, h)
				rgba.Set(xm, ym, colorer(p))
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

func drawLineBetween(x1, y1, x2, y2 int, colorer Colorer, thickness int) *image.RGBA {

	// Bresenham's line-drawing algorithm from wikipedia
	xDelta := math.Abs(float64(x2 - x1))
	yDelta := math.Abs(float64(y2 - y1))

	if xDelta == 0 && yDelta == 0 {
		rect := image.Rect(0, 0, 1, 1)
		rgba := image.NewRGBA(rect)
		rgba.Set(0, 0, colorer(1.0))
		return rgba
	} else if xDelta == 0 {
		rect := image.Rect(0, 0, 1, int(math.Floor(yDelta)))
		rgba := image.NewRGBA(rect)
		for i := 0.0; i <= math.Floor(yDelta); i++ {
			rgba.Set(0, int(i), colorer(i/yDelta))
		}
		return rgba
	}

	// Todo: document why we add one here
	// It has something to do with zero-height rgbas, but is always useful
	h := int(yDelta) + 1

	rect := image.Rect(0, 0, int(xDelta)+2*thickness, h+2*thickness)
	rgba := image.NewRGBA(rect)

	x2 += thickness
	y2 += thickness
	x1 += thickness
	y1 += thickness

	DrawLineColored(rgba, x1, y1, x2, y2, thickness, colorer)

	return rgba
}
