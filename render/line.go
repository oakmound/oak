package render

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
	"image/draw"
	"math"
)

type Line struct {
	x, y  float64
	r     *image.RGBA
	layer int
}

func NewLine(x1, y1, x2, y2 float64, c color.Color) *Line {

	xDelta := math.Abs(x2 - x1)
	rgba := drawLineBetween(int(y1), int(xDelta), int(y2), c)

	return &Line{
		math.Max(x1, x2),
		math.Max(y1, y2),
		rgba,
		0,
	}
}

func (ln *Line) GetRGBA() *image.RGBA {
	return ln.r
}

func (ln *Line) Draw(buff screen.Buffer) {
	draw.Draw(buff.RGBA(), buff.Bounds(),
		ln.r, image.Point{int(ln.x),
			int(ln.y)}, draw.Over)
}

func (ln *Line) ShiftX(x float64) {
	ln.x += x
}
func (ln *Line) ShiftY(y float64) {
	ln.y += y
}

func (ln *Line) GetLayer() int {
	return ln.layer
}

func (ln *Line) SetLayer(l int) {
	ln.layer = l
}

func (ln *Line) UnDraw() {
	ln.layer = -1
}

func (ln *Line) SetPos(x, y float64) {
	ln.x = x
	ln.y = y
}

// x1 is always 0
// either y1 or y2 is always 0
func drawLineBetween(y1, x2, y2 int, c color.Color) *image.RGBA {

	// Bresenham's line-drawing algorithm from wikipedia
	xDelta := float64(x2)
	yDelta := math.Abs(float64(y2 - y1))

	if xDelta == 0 && yDelta == 0 {
		rect := image.Rect(0, 0, 1, 1)
		rgba := image.NewRGBA(rect)
		rgba.Set(0, 0, c)
		return rgba
	}

	xSlope := -1
	if x2 < 0 {
		xSlope = 1
	}
	ySlope := -1
	h := int(yDelta) + 1
	if y2 == 0 {
		ySlope = 1
	}
	rect := image.Rect(0, 0, x2, h)
	rgba := image.NewRGBA(rect)

	err := xDelta - yDelta
	var err2 float64
	y3 := y2
	for i := 0; true; i++ {
		dlog.Verb("Setting ", x2, " ", y2-h)
		rgba.Set(x2, y3-y2, c)
		if x2 == 0 && y2 == y1 {
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

	dlog.Verb(y1, x2, y2, rgba)

	return rgba
}
