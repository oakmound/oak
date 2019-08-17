package render

import (
	"github.com/akavel/polyclip-go"
	"github.com/oakmound/oak/alg/floatgeom"
)

// A DrawPolygon is used to determine whether elements should be drawn, defining
// a polygonal area for what things should be visible.
type DrawPolygon struct {
	usingDrawPolygon bool
	drawPolygon      polyclip.Polygon
}

// SetDrawPolygon sets the draw polygon and flags that draw functions
// should check for containment in the polygon before drawing elements
func (dp *DrawPolygon) SetDrawPolygon(p polyclip.Polygon) {
	dp.usingDrawPolygon = true
	dp.drawPolygon = p
}

// ClearDrawPolygon will stop checking the set draw polygon for whether elements
// should be drawn to screen. If SetDrawPolygon was not called before this was
// called, this does nothing.
// This may in the future be called at the start of new scenes.
func (dp *DrawPolygon) ClearDrawPolygon() {
	dp.usingDrawPolygon = false
}

// DrawPolygonDim returns the dimensions of this draw polygon, or (0,0)->(0,0)
// if there is no draw polygon in use.
func (dp *DrawPolygon) DrawPolygonDim() floatgeom.Rect2 {
	if !dp.usingDrawPolygon {
		return floatgeom.Rect2{}
	}
	mbr := dp.drawPolygon.BoundingBox()
	return floatgeom.NewRect2(mbr.Min.X, mbr.Min.Y, mbr.Max.X, mbr.Max.Y)
}

// InDrawPolygon returns whehter a coordinate and dimension set should be drawn
// given the draw polygon
func (dp *DrawPolygon) InDrawPolygon(xi, yi, x2i, y2i int) bool {
	if dp.usingDrawPolygon {
		x := float64(xi)
		y := float64(yi)
		x2 := float64(x2i)
		y2 := float64(y2i)
		p2 := polyclip.Polygon{{{X: x, Y: y}, {X: x, Y: y2}, {X: x2, Y: y2}, {X: x2, Y: y}}}
		intsct := dp.drawPolygon.Construct(polyclip.INTERSECTION, p2)
		return len(intsct) != 0
	}
	return true
}
