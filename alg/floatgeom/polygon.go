package floatgeom

import (
	"github.com/oakmound/oak/v4/alg"
)

// A Polygon2 is a series of points in 2D space.
type Polygon2 struct {
	// Bounding is a cached bounding box calculated from the input points
	// It is exported for convenience, but should be modified with care
	Bounding Rect2
	// The component points of the polygon. If modified, Bounding should
	// be updated with NewBoundingRect2.
	Points      []Point2
	rectangular bool
}

// NewPolygon2 is a helper method to construct a valid polygon. Polygons
// cannot contain less than 3 points.
func NewPolygon2(p1, p2, p3 Point2, pn ...Point2) Polygon2 {
	pts := append([]Point2{p1, p2, p3}, pn...)
	bounding := NewBoundingRect2(pts...)
	return Polygon2{
		Bounding:    bounding,
		Points:      pts,
		rectangular: isRectangular(pts...),
	}
}

// Contains returns whether or not the current Polygon contains the passed in Point.
// If it is known that the polygon is convex, ConvexContains should be preferred for
// performance.
func (pg Polygon2) Contains(x, y float64) (contains bool) {
	if !pg.Bounding.Contains(Point2{x, y}) {
		return
	}

	j := len(pg.Points) - 1
	for i := 0; i < len(pg.Points); i++ {
		tp1 := pg.Points[i]
		tp2 := pg.Points[j]
		if (tp1.Y() > y) != (tp2.Y() > y) { // Three comparisons
			if x < (tp2.X()-tp1.X())*(y-tp1.Y())/(tp2.Y()-tp1.Y())+tp1.X() { // One Comparison, Four add/sub, Two mult/div
				contains = !contains
			}
		}
		j = i
	}
	return
}

// ConvexContains returns whether the given point is contained by the input polygon.
// It assumes the polygon is convex.
func (pg Polygon2) ConvexContains(x, y float64) bool {

	p := Point2{x, y}

	if !pg.Bounding.Contains(p) {
		return false
	}

	prev := 0
	for i := 0; i < len(pg.Points); i++ {
		tp1 := pg.Points[i]
		tp2 := pg.Points[(i+1)%len(pg.Points)]
		tp3 := tp2.Sub(tp1)
		tp4 := p.Sub(tp1)
		cur := getSide(tp3, tp4)
		if cur == 0 {
			return false
		} else if prev == 0 {
		} else if prev != cur {
			return false
		}
		prev = cur
	}
	return true
}

// TODO: rename this to its real math name, export it
func getSide(a, b Point2) int {
	x := a.X()*b.Y() - a.Y()*b.X()
	if x == 0 {
		return 0
	} else if x < 1 {
		return -1
	} else {
		return 1
	}
}

// OverlappingRectCollides returns whether a Rect2 intersects or is contained by this Polygon.
// This method differs from RectCollides because it assumes that we already know r overlaps with pg.Bounding.
// It is only valid for convex polygons.
func (pg Polygon2) OverlappingRectCollides(r Rect2) bool {
	if pg.rectangular {
		return true
	}
	diags := [][2]Point2{
		{
			{r.Min.X(), r.Max.Y()},
			{r.Max.X(), r.Min.Y()},
		}, {
			r.Min,
			r.Max,
		},
	}
	last := pg.Points[len(pg.Points)-1]
	for i := 0; i < len(pg.Points); i++ {
		next := pg.Points[i]
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

// RectCollides returns whether a Rect2 intersects or is contained by this Polygon.
// It is only valid for convex polygons.
func (pg Polygon2) RectCollides(r Rect2) bool {
	x := float64(r.Min.X())
	y := float64(r.Min.Y())
	x2 := float64(r.Max.X())
	y2 := float64(r.Max.Y())

	dx := pg.Bounding.Min.X()
	dy := pg.Bounding.Min.Y()
	dx2 := pg.Bounding.Max.X()
	dy2 := pg.Bounding.Max.Y()

	overlapX := false
	if x > dx {
		if x < dx2 {
			overlapX = true
		}
	} else {
		if dx < x2 {
			overlapX = true
		}
	}
	if !overlapX {
		return false
	}
	if y > dy {
		if y < dy2 {
			return pg.OverlappingRectCollides(r)
		}
	} else {
		if dy < y2 {
			return pg.OverlappingRectCollides(r)
		}
	}
	return false
}

func isRectangular(pts ...Point2) bool {
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

func orient(p1, p2, p3 Point2) int8 {
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
