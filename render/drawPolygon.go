package render

import (
	clip "github.com/akavel/polyclip-go"
)

var (
	usingDrawPolygon = false
	drawPolygon      clip.Polygon
)

// SetDrawPolygon sets the draw polygon and flags that draw functions
// should check for containment in the polygon before running
func SetDrawPolygon(p clip.Polygon) {
	usingDrawPolygon = true
	drawPolygon = p
}

// DrawPolygonDim returns the dimensions of the draw polygon, or all zeroes
// if there is none.
func DrawPolygonDim() (int, int, int, int) {
	if !usingDrawPolygon {
		return 0, 0, 0, 0
	}
	mbr := drawPolygon.BoundingBox()
	return int(mbr.Min.X), int(mbr.Min.Y), int(mbr.Max.X), int(mbr.Max.Y)
}

// InDrawPolygon returns whehter a coordinate and dimension set should be drawn
// given the draw polygon
func InDrawPolygon(xi, yi, x2i, y2i int) bool {
	if usingDrawPolygon {
		x := float64(xi)
		y := float64(yi)
		x2 := float64(x2i)
		y2 := float64(y2i)
		p2 := clip.Polygon{{{X: x, Y: y}, {X: x, Y: y2}, {X: x2, Y: y2}, {X: x2, Y: y}}}
		intsct := drawPolygon.Construct(clip.INTERSECTION, p2)
		return len(intsct) != 0
	}
	return true
}
