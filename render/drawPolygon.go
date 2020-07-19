package render

import (
	"github.com/akavel/polyclip-go"
	"github.com/oakmound/oak/v2/alg"
	"github.com/oakmound/oak/v2/alg/floatgeom"
)

// A DrawPolygon is used to determine whether elements should be drawn, defining
// a polygonal area for what things should be visible.
type DrawPolygon struct {
	usingDrawPolygon bool
	drawPolygon      []floatgeom.Point2
	dims             floatgeom.Rect2
	rectangular      bool
}

// SetDrawPolygon sets the draw polygon and flags that draw functions
// should check for containment in the polygon before drawing elements
// Deprecated: use SetPolygon instead
func (dp *DrawPolygon) SetDrawPolygon(p polyclip.Polygon) {
	// get []floatgeom.Point2
	poly := make([]floatgeom.Point2, 0, len(p))
	for _, c := range p {
		for _, pt := range c {
			poly = append(poly, floatgeom.Point2{pt.X, pt.Y})
		}
	}
	dp.SetPolygon(poly)
}

// SetPolygon sets the draw polygon and flags that draw functions
// should check for containment in the polygon before drawing elements.
func (dp *DrawPolygon) SetPolygon(poly []floatgeom.Point2) {
	dp.usingDrawPolygon = true
	dp.dims = floatgeom.NewBoundingRect2(poly...)
	dp.drawPolygon = poly
	dp.rectangular = isRectangular(poly...)
}

func isRectangular(pts ...floatgeom.Point2) bool {
	last := pts[len(pts)-1]
	for _, pt := range pts {
		// The last point needs to share an x or y value with this point
		if !alg.F64eq(pt.X(), last.X()) && !alg.F64eq(pt.Y(), last.Y()) {
			return false
		}
		last = pt
	}
	return true
}

// ClearDrawPolygon will stop checking the set draw polygon for whether elements
// should be drawn to screen. If SetDrawPolygon was not called before this was
// called, this does nothing.
// This may in the future be called at the start of new scenes.
func (dp *DrawPolygon) ClearDrawPolygon() {
	dp.usingDrawPolygon = false
	dp.dims = floatgeom.Rect2{}
	dp.rectangular = false
}

// DrawPolygonDim returns the dimensions of this draw polygon, or (0,0)->(0,0)
// if there is no draw polygon in use. Deprecated: Use DrawPolygonBounds instead
func (dp *DrawPolygon) DrawPolygonDim() floatgeom.Rect2 {
	return dp.dims
}

// DrawPolygonBounds returns the dimensions of this draw polygon, or (0,0)->(0,0)
// if there is no draw polygon in use. 
func (dp *DrawPolygon) DrawPolygonBounds() floatgeom.Rect2 {
	return dp.dims
}

// InDrawPolygon returns whehter a coordinate and dimension set should be drawn
// given the draw polygon
func (dp *DrawPolygon) InDrawPolygon(xi, yi, x2i, y2i int) bool {
	if dp.usingDrawPolygon {
		x := float64(xi)
		y := float64(yi)
		x2 := float64(x2i)
		y2 := float64(y2i)

		dx := dp.dims.Min.X()
		dy := dp.dims.Min.Y()
		dx2 := dp.dims.Max.X()
		dy2 := dp.dims.Max.Y()

		dimOverlap := false
		if x > dx {
			if x < dx2 {
				dimOverlap = true
			}
		} else {
			if dx < x2 {
				dimOverlap = true
			}
		}
		if y > dy {
			if y < dy2 {
				dimOverlap = true
			}
		} else {
			if dy < y2 {
				dimOverlap = true
			}
		}
		if !dimOverlap {
			return false
		}
		if dp.rectangular {
			return true
		}
		r := floatgeom.NewRect2(x, y, x2, y2)
		diags := [][2]floatgeom.Point2{
			{
				{r.Min.X(), r.Max.Y()},
				{r.Max.X(), r.Min.Y()},
			}, {
				r.Min,
				r.Max,
			},
		}
		last := dp.drawPolygon[len(dp.drawPolygon)-1]
		for i := 0; i < len(dp.drawPolygon); i++ {
			next := dp.drawPolygon[i]
			if r.Contains(next) {
				return true
			}
			// Checking line segment from last to next
			for _, diag := range diags {
				if orient(diag[0], diag[1], last) != orient(diag[0], diag[1], next) &&
					orient(next, last, diag[0]) != orient(next, last, diag[1]) {
					return true
				}
			}

			last = next
		}
		return false
	}
	return true
}

func orient(p1, p2, p3 floatgeom.Point2) int8 {
	val := (p2.Y()-p1.Y())*(p3.X()-p2.X()) -
		(p2.X()-p1.X())*(p3.Y()-p2.Y())
	switch {
	case val < 0:
		return 2
	case val > 0:
		return 1
	default:
		return 0
	}
}
