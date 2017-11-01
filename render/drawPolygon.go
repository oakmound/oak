package render

import (
	"github.com/akavel/polyclip-go"
)

var (
	usingDrawPolygon = false
	drawPolygon      polyclip.Polygon
)

// Todo 2.0: draw polygons should be stack or stack-item specific, not global

// SetDrawPolygon sets the draw polygon and flags that draw functions
// should check for containment in the polygon before drawing elements
func SetDrawPolygon(p polyclip.Polygon) {
	usingDrawPolygon = true
	drawPolygon = p
}

// ClearDrawPolygon will stop checking the set draw polygon for whether elements
// should be drawn to screen. If SetDrawPolygon was not called before this was
// called, this does nothing.
// This may in the future be called at the start of new scenes.
func ClearDrawPolygon() {
	usingDrawPolygon = false
}

// DrawPolygonDim returns the dimensions of the draw polygon, or all zeroes
// if there is none.
// Todo 2.0: This should return a rectangle instead of four elements.
func DrawPolygonDim() (float64, float64, float64, float64) {
	if !usingDrawPolygon {
		return 0, 0, 0, 0
	}
	mbr := drawPolygon.BoundingBox()
	return mbr.Min.X, mbr.Min.Y, mbr.Max.X, mbr.Max.Y
}

// InDrawPolygon returns whehter a coordinate and dimension set should be drawn
// given the draw polygon
func InDrawPolygon(xi, yi, x2i, y2i int) bool {
	if usingDrawPolygon {
		x := float64(xi)
		y := float64(yi)
		x2 := float64(x2i)
		y2 := float64(y2i)
		p2 := polyclip.Polygon{{{X: x, Y: y}, {X: x, Y: y2}, {X: x2, Y: y2}, {X: x2, Y: y}}}
		intsct := drawPolygon.Construct(polyclip.INTERSECTION, p2)
		return len(intsct) != 0
	}
	return true
}
